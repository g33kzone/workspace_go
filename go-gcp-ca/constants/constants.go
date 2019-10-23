package constants

import "time"

// DeliverProjectID is GCP Project for odj-deliver-srv
const DeliverProjectID = "devops-deliver"

// LiveSVPCProjectID is GCP Project for Vulnerabilty Scanning Pub/Sub
const LiveSVPCProjectID = "sit-odj-live-svpc-host"

// TopicOccurrence - Occurrence Topic name
const TopicOccurrence = "container-analysis-occurrences-v1beta1"

// SubscriptionOccurrence - Occurrence Subscription name
const SubscriptionOccurrence = "OccurrencesSubscription"

//SQLNoLockableBuilds log message
const SQLNoLockableBuilds = "No free rows available in build for locking"

//SQLNoBuildStageRecords no build stage records log message
const SQLNoBuildStageRecords = "No build stage records available in build stage details"

//WaitBeforeResume is Time for which monitor job sleeps before checking again if CI/CD is paused
const WaitBeforeResume = 15 * time.Second

//OdjCIPause is flag in environment variable to check if CI is paused
const OdjCIPause = "ODJ_CI_PAUSE"

//OdjCDPause is flag in environment variable to check if CI is paused
const OdjCDPause = "ODJ_CD_PAUSE"

//Criticality to be sent for fetching Project ID from GCP
const Criticality = "work"

const BitBucketUser = "g33kzone"

const URL = "https://example.com/"

const WebHookEvent = "repo:push"

const AuthUser = "manish6385@gmail.com"

const AuthPwd = "welcome123"
