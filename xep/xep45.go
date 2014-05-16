// XEP-0045: Multi-User Chat
// http://xmpp.org/extensions/xep-0045.html
package xep

import (
	"encoding/xml"
	"log"
	"strconv"
)

type MUCX struct {
	XMLName xml.Name `xml:"http://jabber.org/protocol/muc x"`
}

func (_ MUCX) Name() string {
	return "x"
}

func (_ MUCX) FullName() string {
	return "http://jabber.org/protocol/muc x"
}

type ChatRoom struct {
	Jid       string
	Name      string
	Features  []string
	Occupants []string
	Info      *RoomInfo
	Config    *RoomConfig
}

func NewChatRoom(jid, name string) *ChatRoom {
	return &ChatRoom{
		Jid:  jid,
		Name: name,
	}
}

type RoomInfo struct {
	Description   string   // Short Description of Room
	Subject       string   // Current Discussion Topic
	ChangeSubject bool     // The room subject can be modified by participants
	ContactJid    []string // Contact Addresses (normally, room owner or owners)
	MaxHistory    int      // Maximum Number of History Messages Returned by Room
	Lang          string   // Natural Language for Room Discussions
	LdapGroup     string   // An associated LDAP group that defines room membership
	Logs          string   // URL for Archived Discussion Logs
	Occupants     int      // Current Number of Occupants in Room
	CreateTime    string   // Creation date
}

func ParseRoomInfo(form *XFormData) *RoomInfo {
	info := &RoomInfo{}
	if form == nil {
		return info
	}

	for _, field := range form.Fields {
		switch field.Var {
		case "FORM_TYPE":
		case "muc#maxhistoryfetch":
			info.MaxHistory, _ = strconv.Atoi(field.Value[0])
		case "muc#roominfo_contactjid":
			info.ContactJid = append(info.ContactJid, field.Value...)
		case "muc#roominfo_description":
			info.Description = field.Value[0]
		case "muc#roominfo_lang":
			info.Lang = field.Value[0]
		case "muc#roominfo_ldapgroup":
			info.LdapGroup = field.Value[0]
		case "muc#roominfo_logs":
			info.Logs = field.Value[0]
		case "muc#roominfo_occupants":
			info.Occupants, _ = strconv.Atoi(field.Value[0])
		case "muc#roominfo_subject":
			info.Subject = field.Value[0]
		case "muc#roominfo_changesubject", "muc#roominfo_subjectmod":
			info.ChangeSubject, _ = strconv.ParseBool(field.Value[0])
		case "x-muc#roominfo_creationdate":
			info.CreateTime = field.Value[0]
		default:
			log.Println("Unknown muc roominfo:", field.Var)
		}
	}

	return info
}

type RoomConfig struct {
	MaxHistory      int      // Maximum Number of History Messages Returned by Room
	AllowPM         bool     // Roles that May Send Private Messages
	AllowInvites    bool     // Whether to Allow Occupants to Invite Others
	ChangeSubject   bool     // Whether to Allow Occupants to Change Subject
	EnableLogging   bool     // Whether to Enable Public Logging of Room Conversations
	GetMemberList   []string // Roles and Affiliations that May Retrieve Member List
	Lang            string   // Natural Language for Room Discussions
	Pubsub          string   // XMPP URI of Associated Publish-Subcribe Node
	MaxUsers        int      // Maximum Number of Room Occupants
	MembersOnly     bool     // Whether to Make Room Members-Only
	Moderated       bool     // Whether to Make Room Moderated
	PassRequired    bool     // Whether a Password is Required to Enter
	Persistent      bool     // Whether to Make Room Persistent
	Broadcast       []string // Roles for which Presence is Broadcasted
	PublicSearching bool     // Whether to Allow Public Searching for Room
	Admins          []string // Full List of Room Admins
	Description     string   // Short Description of Room
	Name            string   // Natural-Language Room Name
	Owners          []string // Full List of Room Owners
	Password        string   // The Room Password
	Whois           string   // Affiliations that May Discover Real JIDs of Occupants
}
