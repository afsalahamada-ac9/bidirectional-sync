package Entity

// EventOperation represents a single operation on an event
type EventOperation struct {
	Operation *string      `json:"operation"`
	Value     EventDetails `json:"value"`
}

// represents the details of an event
type EventDetails struct {
	AwsExternalID             *string `json:"Aws_External_Id__c"`
	MaxAttendees              int     `json:"Max_Attendees__c"`
	RegistrationStartDateTime *string `json:"Registration_Start_Date_Time__c"`
	RegistrationEndDateTime   *string `json:"Registration_End_Date_Time__c"` // Nullable
	Location                  *string `json:"Location__c"`
	EventEndDate              *string `json:"Event_End_Date__c"`
	EventStartDate            *string `json:"Event_Start_Date__c"`
	Notes                     *string `json:"Notes__c"`
	Status                    *string `json:"Status__c"`
	EventStartTime            *string `json:"Event_Start_Time__c"`
	EventEndTime              *string `json:"Event_End_Time__c"`
	Timezone                  *string `json:"Timezone__c"`
	ContactPerson             *string `json:"Contact_Person__c"`
	Organizer                 *string `json:"Organizer__c"`
	EventStartDateTimeGMT     *string `json:"Event_Start_Date_Time_GMT__c"`
	EventEndDateTimeGMT       *string `json:"Event_End_Date_Time_GMT__c"`
}

// the main JSON object
type SFData struct {
	Object *string          `json:"object"`
	Value  []EventOperation `json:"value"`
}
