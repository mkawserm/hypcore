package mcodes

/*
Integer 1 - 100 is reserved for Internal Hyper Core Error group
*/

const AuthGroupCode = 1
const AuthQueryMustBeUsingPostRequest = 1
const AuthRequestBodyReadError = 2
const AuthRequestBodyParseError = 3
const AuthGraphQLExecutionError = 4

const GraphQLGroupCode = 2
const GraphQLNoAuthorizationHeaderFound = 1
const GraphQLNoUIDFoundFromToken = 2
const GraphQLInvalidAuthorizationData = 3
const GraphQLQueryMustBeUsingPostRequest = 4
const GraphQLRequestBodyReadError = 5
const GraphQLRequestBodyParseError = 6
const GraphQLExecutionError = 7

const GraphQLWSGroupCode = 3
const GraphQLWSMessageParseError = 1
const GraphQLWSExecutionError = 2
