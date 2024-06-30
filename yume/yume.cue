if #Meta.Environment.Type == "development" {
    Prod: false
    LogPotentialBug: true
}
if #Meta.Environment.Type == "production" {
    Prod: true
    LogPotentialBug: false
}
if #Meta.Environment.Type == "ephemeral" {
    Prod: false
    LogPotentialBug: true
}
