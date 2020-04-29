package protocol

import (
	"fmt"
	"strings"
)

// https://dev.mysql.com/doc/internals/en/capability-flags.html

type CapabilityFlag uint32

/**
	Each flag is just an number, that can be represented just by having a single bit ON.
	It allows as to use fast bitwise operations. Each flag is just a number with applied << operator,
	that is equivalent of multiply by 2

	1 = 00000001
	2 = 00000010
	4 = 00000100
	...

	To check if the flag is set, we use & operator

	00000111 & 00000001 = 1 => true
	00000111 & 01000000 = 0 => false
 */

func (r CapabilityFlag) Has(flag CapabilityFlag) bool {
	return r & flag != 0
}

func (r CapabilityFlag) String() string {
	var names []string

	for i := uint64(1); i <= uint64(1) << 31; i = i << 1 {
		name, ok := flags[CapabilityFlag(i)]; if ok {
			names = append(names, fmt.Sprintf("0x%08x - %032b - %s", i, i, name))
		}
	}

	return strings.Join(names, "\n")
}

const (
	clientLongPassword CapabilityFlag = 1 << iota
	clientFoundRows
	clientLongFlag
	clientConnectWithDB
	clientNoSchema
	clientCompress
	clientODBC
	clientLocalFiles
	clientIgnoreSpace
	clientProtocol41
	clientInteractive
	clientSSL
	clientIgnoreSIGPIPE
	clientTransactions
	clientReserved
	clientSecureConn
	clientMultiStatements
	clientMultiResults
	clientPSMultiResults
	clientPluginAuth
	clientConnectAttrs
	clientPluginAuthLenEncClientData
	clientCanHandleExpiredPasswords
	clientSessionTrack
	clientDeprecateEOF
)

var flags = map[CapabilityFlag]string{
	clientLongPassword: "clientLongPassword",
	clientFoundRows: "clientFoundRows",
	clientLongFlag: "clientLongFlag",
	clientConnectWithDB: "clientConnectWithDB",
	clientNoSchema: "clientNoSchema",
	clientCompress: "clientCompress",
	clientODBC: "clientODBC",
	clientLocalFiles: "clientLocalFiles",
	clientIgnoreSpace: "clientIgnoreSpace",
	clientProtocol41: "clientProtocol41",
	clientInteractive: "clientInteractive",
	clientSSL: "clientSSL",
	clientIgnoreSIGPIPE: "clientIgnoreSIGPIPE",
	clientTransactions: "clientTransactions",
	clientReserved: "clientReserved",
	clientSecureConn: "clientSecureConn",
	clientMultiStatements: "clientMultiStatements",
	clientMultiResults: "clientMultiResults",
	clientPSMultiResults: "clientPSMultiResults",
	clientPluginAuth: "clientPluginAuth",
	clientConnectAttrs: "clientConnectAttrs",
	clientPluginAuthLenEncClientData: "clientPluginAuthLenEncClientData",
	clientCanHandleExpiredPasswords: "clientCanHandleExpiredPasswords",
	clientSessionTrack: "clientSessionTrack",
	clientDeprecateEOF: "clientDeprecateEOF",
}