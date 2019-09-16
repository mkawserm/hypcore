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

const GraphQLWSUpgradeGroupCode = 4
const GraphQLWSUpgradeRequestMethodError = 1
const GraphQLWSUpgradeBadProtocol = 2
const GraphQLWSUpgradeHostNotFound = 3
const GraphQLWSUpgradeNoUpgradeHeaderFound = 4
const GraphQLWSUpgradeNoConnectionHeaderFound = 5
const GraphQLWSUpgradeNoAuthorizationHeaderFound = 6
const GraphQLWSUpgradeNoUIDFoundFromToken = 7
const GraphQLWSUpgradeInvalidAuthorizationData = 8
const GraphQLWSUpgradePathNotFound = 9

const DynamicViewGroupCode = 5
const DynamicViewNoAuthorizationHeaderFound = 1
const DynamicViewNoUIDFoundFromToken = 2
const DynamicViewInvalidAuthorizationData = 3
const DynamicViewPathNotFound = 4
