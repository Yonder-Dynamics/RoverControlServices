package util

//Serializable describes an object which can convert itself to and from
//an array of bytes
type Serializable interface {
	//Serialize returns a byte-array representation of the object
	Serialize() []byte

	//Deserialize transforms the object using its byte-array representation
	Deserialize([]byte)
}
