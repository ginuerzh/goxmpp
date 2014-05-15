// XEP-0045: Multi-User Chat
// http://xmpp.org/extensions/xep-0045.html
package xep

import ()

type Room struct {
	Jid      string
	Name     string
	Features []string
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
