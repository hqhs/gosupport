package main

import (
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"time"
)

const (
	bcryptCost = 10
	// Default store names for database
	// e.g. tables names for postgres, collection names for mongo etc
	helpdeskerStoreName   = "helpdeskers"
	telegramUserStoreName = "telegramusers"
	ticketStoreName       = "tickets"
	messageStoreName      = "messages"
	botStoreName          = "bots"
)

var (
	// Database TODO
	Database *mgo.Database
	// Session is persistent connection to mongo database
	Session *mgo.Session
	mongo   *mgo.DialInfo
)

// InitDatabase parses mongoDBUrl from global Conf (defined in main.go), and sets
// global Session and Mongo variables.
func init() {
	// NOTE this init method is calles before Conf is parsed, so I cant use
	// GetHelpdeskerCollection here
	var err error
	Session, err = mgo.Dial(Conf.MongoDBUrl)
	if err != nil {
		log.Printf("Couldn't connect to database: %v\n", err)
		panic(err)
	}
}

// TelegramUser represents everyone, who uses `/start` command in any bot chat
type TelegramUser struct {
	ID string
}

// Helpdesker TODO: add description
type Helpdesker struct {
	// Django model fields:
	// telegram_name = models.CharField(max_length=140, default='')
	// chat_id = models.CharField(max_length=300, default=0)
	// telegram_id = models.IntegerField(unique=True, default=0)
	// reset_password_token = models.CharField(max_length=100, default='')
	// reset_password_sent_at = models.DateTimeField(default=datetime(1970, 1, 1, tzinfo=pytz.utc))
	// current_sign_in_at = models.DateTimeField(default=datetime(1970, 1, 1, tzinfo=pytz.utc))
	// last_sign_in_at = models.DateTimeField(default=datetime(1970, 1, 1, tzinfo=pytz.utc))
	ID             bson.ObjectId `json:"_id,omitempty" bson:"_id,omitempty"`
	Email          string        `json:"email" binding:"required" bson:"email"`
	Name           string        `json:"name" bson:"name"`
	PasswordHash   string        `json:"password" binding:"required" bson:"password"`
	CreatedAt      time.Time
	UpdatedAt      time.Time
	IsActive       bool
	IsAdmin        bool
	SignIntCount   int
	FailedAttempts int
	// list of bot tokens which user is available to query
	// NOTE what if IsAdmin == true, user could get any bot by default
	AvailableBots   []string
	CurrentSignInIP string
	LastSignInIP    string
}

// GetHelpdeskerCollection return default collection for helpdekser struct
func GetHelpdeskerCollection(session *mgo.Session) *mgo.Collection {
	return session.DB(Conf.DBName).C("helpdeskers")
}

// FetchHelpdesker queries database and return helpdesker instance
func FetchHelpdesker(db *mgo.Database, email string) (Helpdesker, error) {
	collection := db.C(helpdeskerStoreName)
	helpdesker := Helpdesker{}
	err := collection.Find(bson.M{"email": email}).One(&helpdesker)
	return helpdesker, err
}

// SetPassword TODO
func (h *Helpdesker) SetPassword(password string) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcryptCost)
	if err != nil {
		return err
	}
	h.PasswordHash = string(passwordHash)
	return nil
}

// Save TODO
func (h Helpdesker) Save(collection *mgo.Collection) error {
	err := collection.Insert(&h)
	return err
}

// Ticket TODO
type Ticket struct {
}

// Message TODO
type Message struct {
}

// Bot TODO
type Bot struct {
}
