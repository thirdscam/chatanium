package Internal

// Key-Value Local Storage for User
//
// This is a simple key-value storage for user. Modules can use this to store data.
// or you can store key-value data from KV.
//
// Example: /kv <Key> <Value>
type KV struct {
	userID string
}

// Initalize KV
//
// This function will be called when discord bot is started.
// and this Module is Internal, so you don't need to call this function.
func (t *KV) Init(UserID string) {
	if t.userID == "" {
		t.userID = UserID
	}
}

// Set Key-Value for user
func (t *KV) SetKV(Key string, Value string) {
}

// Get all keys for user
func (t *KV) GetAllKeys() {
}

// Get value from key
func (t *KV) GetValue(Key string) {
}

// Get value from key. This function is for third-party modules.
func GetValue(UserID string, Key string) {
	// TODO: Separate Internal Module for third-party module. (Planned)
	// so, should move this function for other outside modules.
}
