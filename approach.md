Converting the operations to Config and convert 
`operations.go - > helpers.go` 

Validate Pre Requisites : 

```go
func (d *LabelExtractionModule) ValidatePrerequisites(ctx context.Context, moduleConfig map[string]interface{}, inputData map[string]interface{}) (map[string]interface{}, error) {

// Validating the tables first , (ID)
// Validate table schema (Their fields & Validate each field & if data is provided validate that)
 ) 
// Validate the chosen operation from the config ( Whether suitable or not , from each copnfig 
}

```

Execute : 

```go
func (d *LabelExtractionModule) Execute(ctx context.Context, moduleExecutionID string, moduleConfig map[string]interface{}, validatedInputData map[string]interface{}) (map[string]interface{}, error) { 
// After validating , Take the operation give it to the function , and and store the output
// Fetching config , Performing the operation , Storing the data frame output) 
}

```

Post Execute :

```go
func (d *LabelExtractionModule) PostExecute(ctx context.Context, moduleExecutionID string, output map[string]interface{}) (datastore.ModuleOutput, error) {
// Return the output , Update in MongoDB for updated data 
// Convert back to JSON to export for the user 
}
```

