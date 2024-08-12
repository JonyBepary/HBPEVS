Block component:
    type BlockHeader struct {
    BlockID uint64
    ElectionID uint64
    Pscode string
    Timestamp int64
    PreHash string
    PostHash string
    ProofOfBlock string
    }

    type CandidateStruct struct {
    	ApplicableCandidates []string
    	Digest               string
    	Signature            string
    	}


    type VoterStruct struct {
    	Voters               []string
    	Pscode               string
    	Seed                 string
    	Digest               string
    	Signature            string
    	}


    type AKS struct {
    	PublicKey            string
    	PrivateKey           string
    	EncryptionAlgorithm  string
    	KeySize              int64
    	SigSize              int64
    	}



    type ElectionCommission struct {
    	PublicKey              string
    	PrivateKey             string
    	Ec_EncryptionAlgorithm string
    	KeySize                int64
    	SigSize                int64
    	Digest                 string
    	}


    type PollingOfficer struct {
    	Name                 string
    	Id                   string
    	PollingStation       string
    	Keys                 *AKS
    	Digest               string
    	Signature            string
    	}

    type ReturningOfficer struct {
    	Name                 string
    	Id                   string
    	RArea                []string
    	Keys                 *AKS
    	Digest               string
    	Signature            string
    	}

    type Command struct {
    	AuthorizeBy          string
    	Action               string
    	Entity               string
    	}

    type Digest struct {
    	Digest               string
    	}

    type ProofOfCommand struct {
    	ProofOfCommand       string
    	}

    type ContractStruct struct {
    	Command              []*Command
    	Digest               []*Digest
    	ProofOfCommand       []*ProofOfCommand
    	}

    type PoBannedList struct {
    	PoHash               []string
    	PoPoa                []string
    	}

    type ElectionTime struct {
    	Start                string
    	Duration             string
    	}

    type SpendedToken struct {
    	Token                []byte
    	Digest               string
    	TokenProof           string
    	}

    type Block struct {
    	Header               *BlockHeader
    	ElectionTime         *ElectionTime
    	CandidateStruct      *CandidateStruct
    	VoterStruct          *VoterStruct
    	ElectionCommission   *ElectionCommission
    	PollingOfficer       []*PollingOfficer
    	ReturningOfficer     []*ReturningOfficer
    	ContractStruct       *ContractStruct
    	PoBannedList         *PoBannedList
    	SpendedToken         *SpendedToken
    	Token                string
    	Vote                 string

    type LiteBlock struct {
    	Header               *BlockHeader
    	ElectionTime         string
    	CandidateStruct      string
    	VoterStruct          string
    	ElectionCommission   string
    	PollingOfficer       []string
    	ReturningOfficer     []string
    	ContractStruct       string
    	PoBannedList         string
    	SpendedToken         string
    	Token                string
    	Vote                 string
    }
    ElectionCommission
    CandidateStruct
    VoterStruct
    PollingOfficer
    ReturningOfficer
    ContractStruct
    PoBannedList
    Vote
    SpendedToken
    Token

Complex structure:
    type BlockChain struct {
        Block []
    }

    type MessageData struct {
    	// shared between all requests
    	ClientVersion string
    	Block         *Block
    	// int64 timestamp = 2;
    	Id                   string
    	Gossip               bool
    	NodeId               string
    	NodePubKey           []byte
    	Sign                 []byte
    	}

    type BlockRequest struct {
    	MessageData *MessageData
    	// method specific data
    	Message              string
    	}

    type BlockResponse struct {
    	MessageData *MessageData
    	// response specific data
    	Message              string
    	}
