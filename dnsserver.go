package main

/*
	Simple DNS Server implemented in Go

	BSD 2-Clause License

	Copyright (c) 2019, Daniel Lorch
	All rights reserved.

	Redistribution and use in source and binary forms, with or without
	modification, are permitted provided that the following conditions are met:

	1. Redistributions of source code must retain the above copyright notice, this
	   list of conditions and the following disclaimer.

	2. Redistributions in binary form must reproduce the above copyright notice,
       this list of conditions and the following disclaimer in the documentation
       and/or other materials provided with the distribution.

	THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
	AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
	IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
	DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
	FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
	DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
	SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
	CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
	OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
	OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

// DNSHeader describes the request/response DNS header
type DNSHeader struct {
	TransactionID  uint16
	Flags          uint16
	NumQuestions   uint16
	NumAnswers     uint16
	NumAuthorities uint16
	NumAdditionals uint16
}

// DNSResourceRecord describes individual records in the request and response of the DNS payload body
type DNSResourceRecord struct {
	DomainName         string
	Type               uint16
	Class              uint16
	TimeToLive         uint32
	ResourceDataLength uint16
	ResourceData       []byte
}

// Type and Class values for DNSResourceRecord
const (
	TypeA                  uint16 = 1 // a host address
	ClassINET              uint16 = 1 // the Internet
	FlagResponse           uint16 = 1 << 15
	UDPMaxMessageSizeBytes uint   = 512 // RFC1035
)

// Pretend to look up values in a database
func dbLookup(queryResourceRecord DNSResourceRecord) ([]DNSResourceRecord, []DNSResourceRecord, []DNSResourceRecord) {
	var answerResourceRecords = make([]DNSResourceRecord, 0)
	var authorityResourceRecords = make([]DNSResourceRecord, 0)
	var additionalResourceRecords = make([]DNSResourceRecord, 0)

	if queryResourceRecord.DomainName == "example.com" && queryResourceRecord.Type == TypeA && queryResourceRecord.Class == ClassINET {
		exampleResourceRecord := DNSResourceRecord{
			DomainName:         "example.com",
			Type:               TypeA,
			Class:              ClassINET,
			TimeToLive:         31337,
			ResourceData:       []byte{3, 1, 3, 7}, // ipv4 address
			ResourceDataLength: 4,
		}

		answerResourceRecords = append(answerResourceRecords, exampleResourceRecord)
	}

	return answerResourceRecords, authorityResourceRecords, additionalResourceRecords
}

// RFC1035: "Domain names in messages are expressed in terms of a sequence
// of labels. Each label is represented as a one octet length field followed
// by that number of octets.  Since every domain name ends with the null label
// of the root, a domain name is terminated by a length byte of zero."
func readDomainName(requestBuffer *bytes.Buffer) (string, error) {
	var domainName string

	b, err := requestBuffer.ReadByte()

	for ; b != 0 && err == nil; b, err = requestBuffer.ReadByte() {
		labelLength := int(b)
		labelBytes := requestBuffer.Next(labelLength)
		labelName := string(labelBytes)

		if len(domainName) == 0 {
			domainName = labelName
		} else {
			domainName += "." + labelName
		}
	}

	return domainName, err
}

// RFC1035: "Domain names in messages are expressed in terms of a sequence
// of labels. Each label is represented as a one octet length field followed
// by that number of octets.  Since every domain name ends with the null label
// of the root, a domain name is terminated by a length byte of zero."
func writeDomainName(responseBuffer *bytes.Buffer, domainName string) error {
	labels := strings.Split(domainName, ".")

	for _, label := range labels {
		labelLength := len(label)
		labelBytes := []byte(label)

		responseBuffer.WriteByte(byte(labelLength))
		responseBuffer.Write(labelBytes)
	}

	err := responseBuffer.WriteByte(byte(0))

	return err
}

func handleDNSClient(requestBytes []byte, serverConn *net.UDPConn, clientAddr *net.UDPAddr) {
	/**
	 * read request
	 */
	var requestBuffer = bytes.NewBuffer(requestBytes)
	var queryHeader DNSHeader
	var queryResourceRecords []DNSResourceRecord

	err := binary.Read(requestBuffer, binary.BigEndian, &queryHeader) // network byte order is big endian

	if err != nil {
		fmt.Println("Error decoding header: ", err.Error())
	}

	queryResourceRecords = make([]DNSResourceRecord, queryHeader.NumQuestions)

	for _, queryResourceRecord := range queryResourceRecords {
		queryResourceRecord.DomainName, err = readDomainName(requestBuffer)

		if err != nil {
			fmt.Println("Error decoding label: ", err.Error())
		}

		queryResourceRecord.Type = binary.BigEndian.Uint16(requestBuffer.Next(2))
		queryResourceRecord.Class = binary.BigEndian.Uint16(requestBuffer.Next(2))
	}

	/**
	 * lookup values
	 */
	var answerResourceRecords = make([]DNSResourceRecord, 0)
	var authorityResourceRecords = make([]DNSResourceRecord, 0)
	var additionalResourceRecords = make([]DNSResourceRecord, 0)

	for _, queryResourceRecord := range queryResourceRecords {
		newAnswerRR, newAuthorityRR, newAdditionalRR := dbLookup(queryResourceRecord)

		answerResourceRecords = append(answerResourceRecords, newAnswerRR...) // three dots cause the two lists to be concatenated
		authorityResourceRecords = append(authorityResourceRecords, newAuthorityRR...)
		additionalResourceRecords = append(additionalResourceRecords, newAdditionalRR...)
	}

	/**
	 * write response
	 */
	var responseBuffer = new(bytes.Buffer)
	var responseHeader DNSHeader

	responseHeader = DNSHeader{
		TransactionID:  queryHeader.TransactionID,
		Flags:          FlagResponse,
		NumQuestions:   queryHeader.NumQuestions,
		NumAnswers:     uint16(len(answerResourceRecords)),
		NumAuthorities: uint16(len(authorityResourceRecords)),
		NumAdditionals: uint16(len(additionalResourceRecords)),
	}

	err = binary.Write(responseBuffer, binary.BigEndian, &responseHeader)

	if err != nil {
		fmt.Println("Error writing to buffer: ", err.Error())
	}

	for _, queryResourceRecord := range queryResourceRecords {
		err = writeDomainName(responseBuffer, queryResourceRecord.DomainName)

		if err != nil {
			fmt.Println("Error writing to buffer: ", err.Error())
		}

		binary.Write(responseBuffer, binary.BigEndian, queryResourceRecord.Type)
		binary.Write(responseBuffer, binary.BigEndian, queryResourceRecord.Class)
	}

	for _, answerResourceRecord := range answerResourceRecords {
		err = writeDomainName(responseBuffer, answerResourceRecord.DomainName)

		if err != nil {
			fmt.Println("Error writing to buffer: ", err.Error())
		}

		binary.Write(responseBuffer, binary.BigEndian, answerResourceRecord.Type)
		binary.Write(responseBuffer, binary.BigEndian, answerResourceRecord.Class)
		binary.Write(responseBuffer, binary.BigEndian, answerResourceRecord.TimeToLive)
		binary.Write(responseBuffer, binary.BigEndian, answerResourceRecord.ResourceDataLength)
		binary.Write(responseBuffer, binary.BigEndian, answerResourceRecord.ResourceData)
	}

	for _, authorityResourceRecord := range authorityResourceRecords {
		err = writeDomainName(responseBuffer, authorityResourceRecord.DomainName)

		if err != nil {
			fmt.Println("Error writing to buffer: ", err.Error())
		}

		binary.Write(responseBuffer, binary.BigEndian, authorityResourceRecord.Type)
		binary.Write(responseBuffer, binary.BigEndian, authorityResourceRecord.Class)
		binary.Write(responseBuffer, binary.BigEndian, authorityResourceRecord.TimeToLive)
		binary.Write(responseBuffer, binary.BigEndian, authorityResourceRecord.ResourceDataLength)
		binary.Write(responseBuffer, binary.BigEndian, authorityResourceRecord.ResourceData)
	}

	for _, additionalResourceRecord := range additionalResourceRecords {
		err = writeDomainName(responseBuffer, additionalResourceRecord.DomainName)

		if err != nil {
			fmt.Println("Error writing to buffer: ", err.Error())
		}

		binary.Write(responseBuffer, binary.BigEndian, additionalResourceRecord.Type)
		binary.Write(responseBuffer, binary.BigEndian, additionalResourceRecord.Class)
		binary.Write(responseBuffer, binary.BigEndian, additionalResourceRecord.TimeToLive)
		binary.Write(responseBuffer, binary.BigEndian, additionalResourceRecord.ResourceDataLength)
		binary.Write(responseBuffer, binary.BigEndian, additionalResourceRecord.ResourceData)
	}

	serverConn.WriteToUDP(responseBuffer.Bytes(), clientAddr)
}

func main() {
	serverAddr, err := net.ResolveUDPAddr("udp", ":1053")

	if err != nil {
		fmt.Println("Error resolving UDP address: ", err.Error())
		os.Exit(1)
	}

	serverConn, err := net.ListenUDP("udp", serverAddr)

	if err != nil {
		fmt.Println("Error listening: ", err.Error())
		os.Exit(1)
	}

	fmt.Println("Listening at: ", serverAddr)

	defer serverConn.Close()

	requestBytes := make([]byte, UDPMaxMessageSizeBytes)

	for {
		_, clientAddr, err := serverConn.ReadFromUDP(requestBytes)

		if err != nil {
			fmt.Println("Error receiving: ", err.Error())
		} else {
			fmt.Println("Received request from ", clientAddr)
			go handleDNSClient(requestBytes, serverConn, clientAddr) // array is value type (call-by-value), i.e. copied
		}
	}
}
