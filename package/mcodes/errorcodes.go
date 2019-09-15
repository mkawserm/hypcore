package mcodes

const AuthGroupCode = 100
const AuthQueryMustBeUsingPostRequest = 1001
const AuthRequestBodyReadError = 1002
const AuthRequestBodyParseError = 1003
const AuthGraphQLExecutionError = 1004

const NoAuthorizationHeaderFound = "HCNAHF400"
const NoUIDFromAuthVerifyInterface = "HCNUFAVI400"
const InvalidAuthorizationData = "HCIAD400"
const InvalidRequestMethod = "HCIRM400"
const FailedToReadRequestBody = "HCFTRRB400"

const WebSocketUpgradeBadRequestMethod = "HCWSUBRM400"
const WebSocketBadProtocol = "HCWSBP400"
const WebSocketNoHostFound = "HCWSNHF400"
const WebSocketNoUpgradeHeaderFound = "HCWSNUHF400"
const WebSocketNoConnectionHeaderFound = "HCWSNCHF400"

const InvalidGraphQLQuery = "HCIGQLQ400"

const HttpNotFound = "HCHNF404"
