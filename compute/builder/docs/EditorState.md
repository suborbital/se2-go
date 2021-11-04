# EditorState

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Lang** | **string** | The language of a runnable | 
**Contents** | **string** | The source code of a runnable | 
**Tests** | Pointer to [**[]TestPayload**](TestPayload.md) | An array of tests | [optional] 

## Methods

### NewEditorState

`func NewEditorState(lang string, contents string, ) *EditorState`

NewEditorState instantiates a new EditorState object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEditorStateWithDefaults

`func NewEditorStateWithDefaults() *EditorState`

NewEditorStateWithDefaults instantiates a new EditorState object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLang

`func (o *EditorState) GetLang() string`

GetLang returns the Lang field if non-nil, zero value otherwise.

### GetLangOk

`func (o *EditorState) GetLangOk() (*string, bool)`

GetLangOk returns a tuple with the Lang field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLang

`func (o *EditorState) SetLang(v string)`

SetLang sets Lang field to given value.


### GetContents

`func (o *EditorState) GetContents() string`

GetContents returns the Contents field if non-nil, zero value otherwise.

### GetContentsOk

`func (o *EditorState) GetContentsOk() (*string, bool)`

GetContentsOk returns a tuple with the Contents field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContents

`func (o *EditorState) SetContents(v string)`

SetContents sets Contents field to given value.


### GetTests

`func (o *EditorState) GetTests() []TestPayload`

GetTests returns the Tests field if non-nil, zero value otherwise.

### GetTestsOk

`func (o *EditorState) GetTestsOk() (*[]TestPayload, bool)`

GetTestsOk returns a tuple with the Tests field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTests

`func (o *EditorState) SetTests(v []TestPayload)`

SetTests sets Tests field to given value.

### HasTests

`func (o *EditorState) HasTests() bool`

HasTests returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


