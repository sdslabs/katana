package flaghandlerservice

const (
	ServiceFail = "Failed to start flag submision service"

	ServiceSuccess = "Flag Submission Service Started at port "

	Connected = "Connected to"

	ClosingError = "Failed to close\n"

	ReadError = "Connection Read error\n"

	InitInstruction = "Connected to Flag Submission Service\nInitiate your session by `init <teamName> <password>`\n"

	InvalidCommand = "Invalid Command\n"

	InvalidParams = "Invalid Login Parameters \n"

	TeamAlreadyExists = "Team is already Logged in\n"

	TeamConnected = "Team successfully connected,\nEnter flags to submit them\n"

	InvalidCreds = "Invalid Team Credentials\n"

	SubmitSuccess = "Flag submitted successfully, points: "

	InvalidFlag = "Invalid FLag\n"

	NoLogin = "Please login Team first \n"

	WriteError = "Failed To Write"

	TotalScore = " Total Score: "
)
