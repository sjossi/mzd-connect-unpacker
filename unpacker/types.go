package unpacker

// CompressionType declares the used compression type. Makes the unpacker aware
// of what tool/algorithm should be used
type CompressionType int

const (
	UNDEFINED CompressionType = iota
	GZIP
)

type FilesystemEntry struct {
	Name     string
	Children []FilesystemEntry
}

// Settings is the struct that represents all the parsed settings
type Settings struct {
	Packageid       int64
	CompressionType CompressionType
	TotalStepsCount int64
}

// Instruction is a parsed instruction line from instructions.ini
type Instruction struct {
	StepNo          int
	InstructionStep InstructionStep
	Arguments       []string
	Steps           int
}

type Instructions struct {
	Count        int
	Instructions []Instruction
}

// InstructionStep are the various steps identified in the instructions.ini
type InstructionStep int

const (
	Execute InstructionStep = iota
	ImageUpdate
	FileUpdate
	BreakPoint
	Copy
	Remove
	Create
	RemoveFolderContent
)

// InstructionSet distinguishes the different types of instruction files
type InstructionSet int

const (
	MainIni InstructionSet = iota
	FilesIni
	ExecuteIni
	BinaryIni
)

type DataStorage struct {
	Count      int
	UPType     string
	SubUPType  string
	ReTransmit string
	NewPackage string
}

// Ini contains all the information found in an ini. Empty values means that
// section was not present in the ini file
type Ini struct {
	RootDir          string
	Filename         string
	Folder           string
	Instructions     Instructions
	Settings         Settings
	Instructions_Ext Instructions
	DataStorage      DataStorage
}
