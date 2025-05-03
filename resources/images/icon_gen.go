package images

import (
	"github.com/ddkwork/golibrary/safemap"
	"golang.org/x/exp/shiny/materialdesign/icons"
)

var IconMap = safemap.NewOrdered[string, []byte](func(yield func(string, []byte) bool) {
	yield("Action3DRotationIcon", icons.Action3DRotation)
	yield("ActionAccessibilityIcon", icons.ActionAccessibility)
	yield("ActionAccessibleIcon", icons.ActionAccessible)
	yield("ActionAccountBalanceIcon", icons.ActionAccountBalance)
	yield("ActionAccountBalanceWalletIcon", icons.ActionAccountBalanceWallet)
	yield("ActionAccountBoxIcon", icons.ActionAccountBox)
	yield("ActionAccountCircleIcon", icons.ActionAccountCircle)
	yield("ActionAddShoppingCartIcon", icons.ActionAddShoppingCart)
	yield("ActionAlarmIcon", icons.ActionAlarm)
	yield("ActionAlarmAddIcon", icons.ActionAlarmAdd)
	yield("ActionAlarmOffIcon", icons.ActionAlarmOff)
	yield("ActionAlarmOnIcon", icons.ActionAlarmOn)
	yield("ActionAllOutIcon", icons.ActionAllOut)
	yield("ActionAndroidIcon", icons.ActionAndroid)
	yield("ActionAnnouncementIcon", icons.ActionAnnouncement)
	yield("ActionAspectRatioIcon", icons.ActionAspectRatio)
	yield("ActionAssessmentIcon", icons.ActionAssessment)
	yield("ActionAssignmentIcon", icons.ActionAssignment)
	yield("ActionAssignmentIndIcon", icons.ActionAssignmentInd)
	yield("ActionAssignmentLateIcon", icons.ActionAssignmentLate)
	yield("ActionAssignmentReturnIcon", icons.ActionAssignmentReturn)
	yield("ActionAssignmentReturnedIcon", icons.ActionAssignmentReturned)
	yield("ActionAssignmentTurnedInIcon", icons.ActionAssignmentTurnedIn)
	yield("ActionAutorenewIcon", icons.ActionAutorenew)
	yield("ActionBackupIcon", icons.ActionBackup)
	yield("ActionBookIcon", icons.ActionBook)
	yield("ActionBookmarkIcon", icons.ActionBookmark)
	yield("ActionBookmarkBorderIcon", icons.ActionBookmarkBorder)
	yield("ActionBugReportIcon", icons.ActionBugReport)
	yield("ActionBuildIcon", icons.ActionBuild)
	yield("ActionCachedIcon", icons.ActionCached)
	yield("ActionCameraEnhanceIcon", icons.ActionCameraEnhance)
	yield("ActionCardGiftcardIcon", icons.ActionCardGiftcard)
	yield("ActionCardMembershipIcon", icons.ActionCardMembership)
	yield("ActionCardTravelIcon", icons.ActionCardTravel)
	yield("ActionChangeHistoryIcon", icons.ActionChangeHistory)
	yield("ActionCheckCircleIcon", icons.ActionCheckCircle)
	yield("ActionChromeReaderModeIcon", icons.ActionChromeReaderMode)
	yield("ActionClassIcon", icons.ActionClass)
	yield("ActionCodeIcon", icons.ActionCode)
	yield("ActionCompareArrowsIcon", icons.ActionCompareArrows)
	yield("ActionCopyrightIcon", icons.ActionCopyright)
	yield("ActionCreditCardIcon", icons.ActionCreditCard)
	yield("ActionDashboardIcon", icons.ActionDashboard)
	yield("ActionDateRangeIcon", icons.ActionDateRange)
	yield("ActionDeleteIcon", icons.ActionDelete)
	yield("ActionDeleteForeverIcon", icons.ActionDeleteForever)
	yield("ActionDescriptionIcon", icons.ActionDescription)
	yield("ActionDNSIcon", icons.ActionDNS)
	yield("ActionDoneIcon", icons.ActionDone)
	yield("ActionDoneAllIcon", icons.ActionDoneAll)
	yield("ActionDonutLargeIcon", icons.ActionDonutLarge)
	yield("ActionDonutSmallIcon", icons.ActionDonutSmall)
	yield("ActionEjectIcon", icons.ActionEject)
	yield("ActionEuroSymbolIcon", icons.ActionEuroSymbol)
	yield("ActionEventIcon", icons.ActionEvent)
	yield("ActionEventSeatIcon", icons.ActionEventSeat)
	yield("ActionExitToAppIcon", icons.ActionExitToApp)
	yield("ActionExploreIcon", icons.ActionExplore)
	yield("ActionExtensionIcon", icons.ActionExtension)
	yield("ActionFaceIcon", icons.ActionFace)
	yield("ActionFavoriteIcon", icons.ActionFavorite)
	yield("ActionFavoriteBorderIcon", icons.ActionFavoriteBorder)
	yield("ActionFeedbackIcon", icons.ActionFeedback)
	yield("ActionFindInPageIcon", icons.ActionFindInPage)
	yield("ActionFindReplaceIcon", icons.ActionFindReplace)
	yield("ActionFingerprintIcon", icons.ActionFingerprint)
	yield("ActionFlightLandIcon", icons.ActionFlightLand)
	yield("ActionFlightTakeoffIcon", icons.ActionFlightTakeoff)
	yield("ActionFlipToBackIcon", icons.ActionFlipToBack)
	yield("ActionFlipToFrontIcon", icons.ActionFlipToFront)
	yield("ActionGTranslateIcon", icons.ActionGTranslate)
	yield("ActionGavelIcon", icons.ActionGavel)
	yield("ActionGetAppIcon", icons.ActionGetApp)
	yield("ActionGIFIcon", icons.ActionGIF)
	yield("ActionGradeIcon", icons.ActionGrade)
	yield("ActionGroupWorkIcon", icons.ActionGroupWork)
	yield("ActionHelpIcon", icons.ActionHelp)
	yield("ActionHelpOutlineIcon", icons.ActionHelpOutline)
	yield("ActionHighlightOffIcon", icons.ActionHighlightOff)
	yield("ActionHistoryIcon", icons.ActionHistory)
	yield("ActionHomeIcon", icons.ActionHome)
	yield("ActionHourglassEmptyIcon", icons.ActionHourglassEmpty)
	yield("ActionHourglassFullIcon", icons.ActionHourglassFull)
	yield("ActionHTTPIcon", icons.ActionHTTP)
	yield("ActionHTTPSIcon", icons.ActionHTTPS)
	yield("ActionImportantDevicesIcon", icons.ActionImportantDevices)
	yield("ActionInfoIcon", icons.ActionInfo)
	yield("ActionInfoOutlineIcon", icons.ActionInfoOutline)
	yield("ActionInputIcon", icons.ActionInput)
	yield("ActionInvertColorsIcon", icons.ActionInvertColors)
	yield("ActionLabelIcon", icons.ActionLabel)
	yield("ActionLabelOutlineIcon", icons.ActionLabelOutline)
	yield("ActionLanguageIcon", icons.ActionLanguage)
	yield("ActionLaunchIcon", icons.ActionLaunch)
	yield("ActionLightbulbOutlineIcon", icons.ActionLightbulbOutline)
	yield("ActionLineStyleIcon", icons.ActionLineStyle)
	yield("ActionLineWeightIcon", icons.ActionLineWeight)
	yield("ActionListIcon", icons.ActionList)
	yield("ActionLockIcon", icons.ActionLock)
	yield("ActionLockOpenIcon", icons.ActionLockOpen)
	yield("ActionLockOutlineIcon", icons.ActionLockOutline)
	yield("ActionLoyaltyIcon", icons.ActionLoyalty)
	yield("ActionMarkUnreadMailboxIcon", icons.ActionMarkUnreadMailbox)
	yield("ActionMotorcycleIcon", icons.ActionMotorcycle)
	yield("ActionNoteAddIcon", icons.ActionNoteAdd)
	yield("ActionOfflinePinIcon", icons.ActionOfflinePin)
	yield("ActionOpacityIcon", icons.ActionOpacity)
	yield("ActionOpenInBrowserIcon", icons.ActionOpenInBrowser)
	yield("ActionOpenInNewIcon", icons.ActionOpenInNew)
	yield("ActionOpenWithIcon", icons.ActionOpenWith)
	yield("ActionPageviewIcon", icons.ActionPageview)
	yield("ActionPanToolIcon", icons.ActionPanTool)
	yield("ActionPaymentIcon", icons.ActionPayment)
	yield("ActionPermCameraMicIcon", icons.ActionPermCameraMic)
	yield("ActionPermContactCalendarIcon", icons.ActionPermContactCalendar)
	yield("ActionPermDataSettingIcon", icons.ActionPermDataSetting)
	yield("ActionPermDeviceInformationIcon", icons.ActionPermDeviceInformation)
	yield("ActionPermIdentityIcon", icons.ActionPermIdentity)
	yield("ActionPermMediaIcon", icons.ActionPermMedia)
	yield("ActionPermPhoneMsgIcon", icons.ActionPermPhoneMsg)
	yield("ActionPermScanWiFiIcon", icons.ActionPermScanWiFi)
	yield("ActionPetsIcon", icons.ActionPets)
	yield("ActionPictureInPictureIcon", icons.ActionPictureInPicture)
	yield("ActionPictureInPictureAltIcon", icons.ActionPictureInPictureAlt)
	yield("ActionPlayForWorkIcon", icons.ActionPlayForWork)
	yield("ActionPolymerIcon", icons.ActionPolymer)
	yield("ActionPowerSettingsNewIcon", icons.ActionPowerSettingsNew)
	yield("ActionPregnantWomanIcon", icons.ActionPregnantWoman)
	yield("ActionPrintIcon", icons.ActionPrint)
	yield("ActionQueryBuilderIcon", icons.ActionQueryBuilder)
	yield("ActionQuestionAnswerIcon", icons.ActionQuestionAnswer)
	yield("ActionReceiptIcon", icons.ActionReceipt)
	yield("ActionRecordVoiceOverIcon", icons.ActionRecordVoiceOver)
	yield("ActionRedeemIcon", icons.ActionRedeem)
	yield("ActionRemoveShoppingCartIcon", icons.ActionRemoveShoppingCart)
	yield("ActionReorderIcon", icons.ActionReorder)
	yield("ActionReportProblemIcon", icons.ActionReportProblem)
	yield("ActionRestoreIcon", icons.ActionRestore)
	yield("ActionRestorePageIcon", icons.ActionRestorePage)
	yield("ActionRoomIcon", icons.ActionRoom)
	yield("ActionRoundedCornerIcon", icons.ActionRoundedCorner)
	yield("ActionRowingIcon", icons.ActionRowing)
	yield("ActionScheduleIcon", icons.ActionSchedule)
	yield("ActionSearchIcon", icons.ActionSearch)
	yield("ActionSettingsIcon", icons.ActionSettings)
	yield("ActionSettingsApplicationsIcon", icons.ActionSettingsApplications)
	yield("ActionSettingsBackupRestoreIcon", icons.ActionSettingsBackupRestore)
	yield("ActionSettingsBluetoothIcon", icons.ActionSettingsBluetooth)
	yield("ActionSettingsBrightnessIcon", icons.ActionSettingsBrightness)
	yield("ActionSettingsCellIcon", icons.ActionSettingsCell)
	yield("ActionSettingsEthernetIcon", icons.ActionSettingsEthernet)
	yield("ActionSettingsInputAntennaIcon", icons.ActionSettingsInputAntenna)
	yield("ActionSettingsInputComponentIcon", icons.ActionSettingsInputComponent)
	yield("ActionSettingsInputCompositeIcon", icons.ActionSettingsInputComposite)
	yield("ActionSettingsInputHDMIIcon", icons.ActionSettingsInputHDMI)
	yield("ActionSettingsInputSVideoIcon", icons.ActionSettingsInputSVideo)
	yield("ActionSettingsOverscanIcon", icons.ActionSettingsOverscan)
	yield("ActionSettingsPhoneIcon", icons.ActionSettingsPhone)
	yield("ActionSettingsPowerIcon", icons.ActionSettingsPower)
	yield("ActionSettingsRemoteIcon", icons.ActionSettingsRemote)
	yield("ActionSettingsVoiceIcon", icons.ActionSettingsVoice)
	yield("ActionShopIcon", icons.ActionShop)
	yield("ActionShopTwoIcon", icons.ActionShopTwo)
	yield("ActionShoppingBasketIcon", icons.ActionShoppingBasket)
	yield("ActionShoppingCartIcon", icons.ActionShoppingCart)
	yield("ActionSpeakerNotesIcon", icons.ActionSpeakerNotes)
	yield("ActionSpeakerNotesOffIcon", icons.ActionSpeakerNotesOff)
	yield("ActionSpellcheckIcon", icons.ActionSpellcheck)
	yield("ActionStarRateIcon", icons.ActionStarRate)
	yield("ActionStarsIcon", icons.ActionStars)
	yield("ActionStoreIcon", icons.ActionStore)
	yield("ActionSubjectIcon", icons.ActionSubject)
	yield("ActionSupervisorAccountIcon", icons.ActionSupervisorAccount)
	yield("ActionSwapHorizIcon", icons.ActionSwapHoriz)
	yield("ActionSwapVertIcon", icons.ActionSwapVert)
	yield("ActionSwapVerticalCircleIcon", icons.ActionSwapVerticalCircle)
	yield("ActionSystemUpdateAltIcon", icons.ActionSystemUpdateAlt)
	yield("ActionTabIcon", icons.ActionTab)
	yield("ActionTabUnselectedIcon", icons.ActionTabUnselected)
	yield("ActionTheatersIcon", icons.ActionTheaters)
	yield("ActionThumbDownIcon", icons.ActionThumbDown)
	yield("ActionThumbUpIcon", icons.ActionThumbUp)
	yield("ActionThumbsUpDownIcon", icons.ActionThumbsUpDown)
	yield("ActionTimelineIcon", icons.ActionTimeline)
	yield("ActionTOCIcon", icons.ActionTOC)
	yield("ActionTodayIcon", icons.ActionToday)
	yield("ActionTollIcon", icons.ActionToll)
	yield("ActionTouchAppIcon", icons.ActionTouchApp)
	yield("ActionTrackChangesIcon", icons.ActionTrackChanges)
	yield("ActionTranslateIcon", icons.ActionTranslate)
	yield("ActionTrendingDownIcon", icons.ActionTrendingDown)
	yield("ActionTrendingFlatIcon", icons.ActionTrendingFlat)
	yield("ActionTrendingUpIcon", icons.ActionTrendingUp)
	yield("ActionTurnedInIcon", icons.ActionTurnedIn)
	yield("ActionTurnedInNotIcon", icons.ActionTurnedInNot)
	yield("ActionUpdateIcon", icons.ActionUpdate)
	yield("ActionVerifiedUserIcon", icons.ActionVerifiedUser)
	yield("ActionViewAgendaIcon", icons.ActionViewAgenda)
	yield("ActionViewArrayIcon", icons.ActionViewArray)
	yield("ActionViewCarouselIcon", icons.ActionViewCarousel)
	yield("ActionViewColumnIcon", icons.ActionViewColumn)
	yield("ActionViewDayIcon", icons.ActionViewDay)
	yield("ActionViewHeadlineIcon", icons.ActionViewHeadline)
	yield("ActionViewListIcon", icons.ActionViewList)
	yield("ActionViewModuleIcon", icons.ActionViewModule)
	yield("ActionViewQuiltIcon", icons.ActionViewQuilt)
	yield("ActionViewStreamIcon", icons.ActionViewStream)
	yield("ActionViewWeekIcon", icons.ActionViewWeek)
	yield("ActionVisibilityIcon", icons.ActionVisibility)
	yield("ActionVisibilityOffIcon", icons.ActionVisibilityOff)
	yield("ActionWatchLaterIcon", icons.ActionWatchLater)
	yield("ActionWorkIcon", icons.ActionWork)
	yield("ActionYoutubeSearchedForIcon", icons.ActionYoutubeSearchedFor)
	yield("ActionZoomInIcon", icons.ActionZoomIn)
	yield("ActionZoomOutIcon", icons.ActionZoomOut)
	yield("AlertAddAlertIcon", icons.AlertAddAlert)
	yield("AlertErrorIcon", icons.AlertError)
	yield("AlertErrorOutlineIcon", icons.AlertErrorOutline)
	yield("AlertWarningIcon", icons.AlertWarning)
	yield("AVAddToQueueIcon", icons.AVAddToQueue)
	yield("AVAirplayIcon", icons.AVAirplay)
	yield("AVAlbumIcon", icons.AVAlbum)
	yield("AVArtTrackIcon", icons.AVArtTrack)
	yield("AVAVTimerIcon", icons.AVAVTimer)
	yield("AVBrandingWatermarkIcon", icons.AVBrandingWatermark)
	yield("AVCallToActionIcon", icons.AVCallToAction)
	yield("AVClosedCaptionIcon", icons.AVClosedCaption)
	yield("AVEqualizerIcon", icons.AVEqualizer)
	yield("AVExplicitIcon", icons.AVExplicit)
	yield("AVFastForwardIcon", icons.AVFastForward)
	yield("AVFastRewindIcon", icons.AVFastRewind)
	yield("AVFeaturedPlayListIcon", icons.AVFeaturedPlayList)
	yield("AVFeaturedVideoIcon", icons.AVFeaturedVideo)
	yield("AVFiberDVRIcon", icons.AVFiberDVR)
	yield("AVFiberManualRecordIcon", icons.AVFiberManualRecord)
	yield("AVFiberNewIcon", icons.AVFiberNew)
	yield("AVFiberPinIcon", icons.AVFiberPin)
	yield("AVFiberSmartRecordIcon", icons.AVFiberSmartRecord)
	yield("AVForward10Icon", icons.AVForward10)
	yield("AVForward30Icon", icons.AVForward30)
	yield("AVForward5Icon", icons.AVForward5)
	yield("AVGamesIcon", icons.AVGames)
	yield("AVHDIcon", icons.AVHD)
	yield("AVHearingIcon", icons.AVHearing)
	yield("AVHighQualityIcon", icons.AVHighQuality)
	yield("AVLibraryAddIcon", icons.AVLibraryAdd)
	yield("AVLibraryBooksIcon", icons.AVLibraryBooks)
	yield("AVLibraryMusicIcon", icons.AVLibraryMusic)
	yield("AVLoopIcon", icons.AVLoop)
	yield("AVMicIcon", icons.AVMic)
	yield("AVMicNoneIcon", icons.AVMicNone)
	yield("AVMicOffIcon", icons.AVMicOff)
	yield("AVMovieIcon", icons.AVMovie)
	yield("AVMusicVideoIcon", icons.AVMusicVideo)
	yield("AVNewReleasesIcon", icons.AVNewReleases)
	yield("AVNotInterestedIcon", icons.AVNotInterested)
	yield("AVNoteIcon", icons.AVNote)
	yield("AVPauseIcon", icons.AVPause)
	yield("AVPauseCircleFilledIcon", icons.AVPauseCircleFilled)
	yield("AVPauseCircleOutlineIcon", icons.AVPauseCircleOutline)
	yield("AVPlayArrowIcon", icons.AVPlayArrow)
	yield("AVPlayCircleFilledIcon", icons.AVPlayCircleFilled)
	yield("AVPlayCircleOutlineIcon", icons.AVPlayCircleOutline)
	yield("AVPlaylistAddIcon", icons.AVPlaylistAdd)
	yield("AVPlaylistAddCheckIcon", icons.AVPlaylistAddCheck)
	yield("AVPlaylistPlayIcon", icons.AVPlaylistPlay)
	yield("AVQueueIcon", icons.AVQueue)
	yield("AVQueueMusicIcon", icons.AVQueueMusic)
	yield("AVQueuePlayNextIcon", icons.AVQueuePlayNext)
	yield("AVRadioIcon", icons.AVRadio)
	yield("AVRecentActorsIcon", icons.AVRecentActors)
	yield("AVRemoveFromQueueIcon", icons.AVRemoveFromQueue)
	yield("AVRepeatIcon", icons.AVRepeat)
	yield("AVRepeatOneIcon", icons.AVRepeatOne)
	yield("AVReplayIcon", icons.AVReplay)
	yield("AVReplay10Icon", icons.AVReplay10)
	yield("AVReplay30Icon", icons.AVReplay30)
	yield("AVReplay5Icon", icons.AVReplay5)
	yield("AVShuffleIcon", icons.AVShuffle)
	yield("AVSkipNextIcon", icons.AVSkipNext)
	yield("AVSkipPreviousIcon", icons.AVSkipPrevious)
	yield("AVSlowMotionVideoIcon", icons.AVSlowMotionVideo)
	yield("AVSnoozeIcon", icons.AVSnooze)
	yield("AVSortByAlphaIcon", icons.AVSortByAlpha)
	yield("AVStopIcon", icons.AVStop)
	yield("AVSubscriptionsIcon", icons.AVSubscriptions)
	yield("AVSubtitlesIcon", icons.AVSubtitles)
	yield("AVSurroundSoundIcon", icons.AVSurroundSound)
	yield("AVVideoCallIcon", icons.AVVideoCall)
	yield("AVVideoLabelIcon", icons.AVVideoLabel)
	yield("AVVideoLibraryIcon", icons.AVVideoLibrary)
	yield("AVVideocamIcon", icons.AVVideocam)
	yield("AVVideocamOffIcon", icons.AVVideocamOff)
	yield("AVVolumeDownIcon", icons.AVVolumeDown)
	yield("AVVolumeMuteIcon", icons.AVVolumeMute)
	yield("AVVolumeOffIcon", icons.AVVolumeOff)
	yield("AVVolumeUpIcon", icons.AVVolumeUp)
	yield("AVWebIcon", icons.AVWeb)
	yield("AVWebAssetIcon", icons.AVWebAsset)
	yield("CommunicationBusinessIcon", icons.CommunicationBusiness)
	yield("CommunicationCallIcon", icons.CommunicationCall)
	yield("CommunicationCallEndIcon", icons.CommunicationCallEnd)
	yield("CommunicationCallMadeIcon", icons.CommunicationCallMade)
	yield("CommunicationCallMergeIcon", icons.CommunicationCallMerge)
	yield("CommunicationCallMissedIcon", icons.CommunicationCallMissed)
	yield("CommunicationCallMissedOutgoingIcon", icons.CommunicationCallMissedOutgoing)
	yield("CommunicationCallReceivedIcon", icons.CommunicationCallReceived)
	yield("CommunicationCallSplitIcon", icons.CommunicationCallSplit)
	yield("CommunicationChatIcon", icons.CommunicationChat)
	yield("CommunicationChatBubbleIcon", icons.CommunicationChatBubble)
	yield("CommunicationChatBubbleOutlineIcon", icons.CommunicationChatBubbleOutline)
	yield("CommunicationClearAllIcon", icons.CommunicationClearAll)
	yield("CommunicationCommentIcon", icons.CommunicationComment)
	yield("CommunicationContactMailIcon", icons.CommunicationContactMail)
	yield("CommunicationContactPhoneIcon", icons.CommunicationContactPhone)
	yield("CommunicationContactsIcon", icons.CommunicationContacts)
	yield("CommunicationDialerSIPIcon", icons.CommunicationDialerSIP)
	yield("CommunicationDialpadIcon", icons.CommunicationDialpad)
	yield("CommunicationEmailIcon", icons.CommunicationEmail)
	yield("CommunicationForumIcon", icons.CommunicationForum)
	yield("CommunicationImportContactsIcon", icons.CommunicationImportContacts)
	yield("CommunicationImportExportIcon", icons.CommunicationImportExport)
	yield("CommunicationInvertColorsOffIcon", icons.CommunicationInvertColorsOff)
	yield("CommunicationLiveHelpIcon", icons.CommunicationLiveHelp)
	yield("CommunicationLocationOffIcon", icons.CommunicationLocationOff)
	yield("CommunicationLocationOnIcon", icons.CommunicationLocationOn)
	yield("CommunicationMailOutlineIcon", icons.CommunicationMailOutline)
	yield("CommunicationMessageIcon", icons.CommunicationMessage)
	yield("CommunicationNoSIMIcon", icons.CommunicationNoSIM)
	yield("CommunicationPhoneIcon", icons.CommunicationPhone)
	yield("CommunicationPhoneLinkEraseIcon", icons.CommunicationPhoneLinkErase)
	yield("CommunicationPhoneLinkLockIcon", icons.CommunicationPhoneLinkLock)
	yield("CommunicationPhoneLinkRingIcon", icons.CommunicationPhoneLinkRing)
	yield("CommunicationPhoneLinkSetupIcon", icons.CommunicationPhoneLinkSetup)
	yield("CommunicationPortableWiFiOffIcon", icons.CommunicationPortableWiFiOff)
	yield("CommunicationPresentToAllIcon", icons.CommunicationPresentToAll)
	yield("CommunicationRingVolumeIcon", icons.CommunicationRingVolume)
	yield("CommunicationRSSFeedIcon", icons.CommunicationRSSFeed)
	yield("CommunicationScreenShareIcon", icons.CommunicationScreenShare)
	yield("CommunicationSpeakerPhoneIcon", icons.CommunicationSpeakerPhone)
	yield("CommunicationStayCurrentLandscapeIcon", icons.CommunicationStayCurrentLandscape)
	yield("CommunicationStayCurrentPortraitIcon", icons.CommunicationStayCurrentPortrait)
	yield("CommunicationStayPrimaryLandscapeIcon", icons.CommunicationStayPrimaryLandscape)
	yield("CommunicationStayPrimaryPortraitIcon", icons.CommunicationStayPrimaryPortrait)
	yield("CommunicationStopScreenShareIcon", icons.CommunicationStopScreenShare)
	yield("CommunicationSwapCallsIcon", icons.CommunicationSwapCalls)
	yield("CommunicationTextSMSIcon", icons.CommunicationTextSMS)
	yield("CommunicationVoicemailIcon", icons.CommunicationVoicemail)
	yield("CommunicationVPNKeyIcon", icons.CommunicationVPNKey)
	yield("ContentAddIcon", icons.ContentAdd)
	yield("ContentAddBoxIcon", icons.ContentAddBox)
	yield("ContentAddCircleIcon", icons.ContentAddCircle)
	yield("ContentAddCircleOutlineIcon", icons.ContentAddCircleOutline)
	yield("ContentArchiveIcon", icons.ContentArchive)
	yield("ContentBackspaceIcon", icons.ContentBackspace)
	yield("ContentBlockIcon", icons.ContentBlock)
	yield("ContentClearIcon", icons.ContentClear)
	yield("ContentContentCopyIcon", icons.ContentContentCopy)
	yield("ContentContentCutIcon", icons.ContentContentCut)
	yield("ContentContentPasteIcon", icons.ContentContentPaste)
	yield("ContentCreateIcon", icons.ContentCreate)
	yield("ContentDeleteSweepIcon", icons.ContentDeleteSweep)
	yield("ContentDraftsIcon", icons.ContentDrafts)
	yield("ContentFilterListIcon", icons.ContentFilterList)
	yield("ContentFlagIcon", icons.ContentFlag)
	yield("ContentFontDownloadIcon", icons.ContentFontDownload)
	yield("ContentForwardIcon", icons.ContentForward)
	yield("ContentGestureIcon", icons.ContentGesture)
	yield("ContentInboxIcon", icons.ContentInbox)
	yield("ContentLinkIcon", icons.ContentLink)
	yield("ContentLowPriorityIcon", icons.ContentLowPriority)
	yield("ContentMailIcon", icons.ContentMail)
	yield("ContentMarkUnreadIcon", icons.ContentMarkUnread)
	yield("ContentMoveToInboxIcon", icons.ContentMoveToInbox)
	yield("ContentNextWeekIcon", icons.ContentNextWeek)
	yield("ContentRedoIcon", icons.ContentRedo)
	yield("ContentRemoveIcon", icons.ContentRemove)
	yield("ContentRemoveCircleIcon", icons.ContentRemoveCircle)
	yield("ContentRemoveCircleOutlineIcon", icons.ContentRemoveCircleOutline)
	yield("ContentReplyIcon", icons.ContentReply)
	yield("ContentReplyAllIcon", icons.ContentReplyAll)
	yield("ContentReportIcon", icons.ContentReport)
	yield("ContentSaveIcon", icons.ContentSave)
	yield("ContentSelectAllIcon", icons.ContentSelectAll)
	yield("ContentSendIcon", icons.ContentSend)
	yield("ContentSortIcon", icons.ContentSort)
	yield("ContentTextFormatIcon", icons.ContentTextFormat)
	yield("ContentUnarchiveIcon", icons.ContentUnarchive)
	yield("ContentUndoIcon", icons.ContentUndo)
	yield("ContentWeekendIcon", icons.ContentWeekend)
	yield("DeviceAccessAlarmIcon", icons.DeviceAccessAlarm)
	yield("DeviceAccessAlarmsIcon", icons.DeviceAccessAlarms)
	yield("DeviceAccessTimeIcon", icons.DeviceAccessTime)
	yield("DeviceAddAlarmIcon", icons.DeviceAddAlarm)
	yield("DeviceAirplaneModeActiveIcon", icons.DeviceAirplaneModeActive)
	yield("DeviceAirplaneModeInactiveIcon", icons.DeviceAirplaneModeInactive)
	yield("DeviceBattery20Icon", icons.DeviceBattery20)
	yield("DeviceBattery30Icon", icons.DeviceBattery30)
	yield("DeviceBattery50Icon", icons.DeviceBattery50)
	yield("DeviceBattery60Icon", icons.DeviceBattery60)
	yield("DeviceBattery80Icon", icons.DeviceBattery80)
	yield("DeviceBattery90Icon", icons.DeviceBattery90)
	yield("DeviceBatteryAlertIcon", icons.DeviceBatteryAlert)
	yield("DeviceBatteryCharging20Icon", icons.DeviceBatteryCharging20)
	yield("DeviceBatteryCharging30Icon", icons.DeviceBatteryCharging30)
	yield("DeviceBatteryCharging50Icon", icons.DeviceBatteryCharging50)
	yield("DeviceBatteryCharging60Icon", icons.DeviceBatteryCharging60)
	yield("DeviceBatteryCharging80Icon", icons.DeviceBatteryCharging80)
	yield("DeviceBatteryCharging90Icon", icons.DeviceBatteryCharging90)
	yield("DeviceBatteryChargingFullIcon", icons.DeviceBatteryChargingFull)
	yield("DeviceBatteryFullIcon", icons.DeviceBatteryFull)
	yield("DeviceBatteryStdIcon", icons.DeviceBatteryStd)
	yield("DeviceBatteryUnknownIcon", icons.DeviceBatteryUnknown)
	yield("DeviceBluetoothIcon", icons.DeviceBluetooth)
	yield("DeviceBluetoothConnectedIcon", icons.DeviceBluetoothConnected)
	yield("DeviceBluetoothDisabledIcon", icons.DeviceBluetoothDisabled)
	yield("DeviceBluetoothSearchingIcon", icons.DeviceBluetoothSearching)
	yield("DeviceBrightnessAutoIcon", icons.DeviceBrightnessAuto)
	yield("DeviceBrightnessHighIcon", icons.DeviceBrightnessHigh)
	yield("DeviceBrightnessLowIcon", icons.DeviceBrightnessLow)
	yield("DeviceBrightnessMediumIcon", icons.DeviceBrightnessMedium)
	yield("DeviceDataUsageIcon", icons.DeviceDataUsage)
	yield("DeviceDeveloperModeIcon", icons.DeviceDeveloperMode)
	yield("DeviceDevicesIcon", icons.DeviceDevices)
	yield("DeviceDVRIcon", icons.DeviceDVR)
	yield("DeviceGPSFixedIcon", icons.DeviceGPSFixed)
	yield("DeviceGPSNotFixedIcon", icons.DeviceGPSNotFixed)
	yield("DeviceGPSOffIcon", icons.DeviceGPSOff)
	yield("DeviceGraphicEqIcon", icons.DeviceGraphicEq)
	yield("DeviceLocationDisabledIcon", icons.DeviceLocationDisabled)
	yield("DeviceLocationSearchingIcon", icons.DeviceLocationSearching)
	yield("DeviceNetworkCellIcon", icons.DeviceNetworkCell)
	yield("DeviceNetworkWiFiIcon", icons.DeviceNetworkWiFi)
	yield("DeviceNFCIcon", icons.DeviceNFC)
	yield("DeviceScreenLockLandscapeIcon", icons.DeviceScreenLockLandscape)
	yield("DeviceScreenLockPortraitIcon", icons.DeviceScreenLockPortrait)
	yield("DeviceScreenLockRotationIcon", icons.DeviceScreenLockRotation)
	yield("DeviceScreenRotationIcon", icons.DeviceScreenRotation)
	yield("DeviceSDStorageIcon", icons.DeviceSDStorage)
	yield("DeviceSettingsSystemDaydreamIcon", icons.DeviceSettingsSystemDaydream)
	yield("DeviceSignalCellular0BarIcon", icons.DeviceSignalCellular0Bar)
	yield("DeviceSignalCellular1BarIcon", icons.DeviceSignalCellular1Bar)
	yield("DeviceSignalCellular2BarIcon", icons.DeviceSignalCellular2Bar)
	yield("DeviceSignalCellular3BarIcon", icons.DeviceSignalCellular3Bar)
	yield("DeviceSignalCellular4BarIcon", icons.DeviceSignalCellular4Bar)
	yield("DeviceSignalCellularConnectedNoInternet0BarIcon", icons.DeviceSignalCellularConnectedNoInternet0Bar)
	yield("DeviceSignalCellularConnectedNoInternet1BarIcon", icons.DeviceSignalCellularConnectedNoInternet1Bar)
	yield("DeviceSignalCellularConnectedNoInternet2BarIcon", icons.DeviceSignalCellularConnectedNoInternet2Bar)
	yield("DeviceSignalCellularConnectedNoInternet3BarIcon", icons.DeviceSignalCellularConnectedNoInternet3Bar)
	yield("DeviceSignalCellularConnectedNoInternet4BarIcon", icons.DeviceSignalCellularConnectedNoInternet4Bar)
	yield("DeviceSignalCellularNoSIMIcon", icons.DeviceSignalCellularNoSIM)
	yield("DeviceSignalCellularNullIcon", icons.DeviceSignalCellularNull)
	yield("DeviceSignalCellularOffIcon", icons.DeviceSignalCellularOff)
	yield("DeviceSignalWiFi0BarIcon", icons.DeviceSignalWiFi0Bar)
	yield("DeviceSignalWiFi1BarIcon", icons.DeviceSignalWiFi1Bar)
	yield("DeviceSignalWiFi1BarLockIcon", icons.DeviceSignalWiFi1BarLock)
	yield("DeviceSignalWiFi2BarIcon", icons.DeviceSignalWiFi2Bar)
	yield("DeviceSignalWiFi2BarLockIcon", icons.DeviceSignalWiFi2BarLock)
	yield("DeviceSignalWiFi3BarIcon", icons.DeviceSignalWiFi3Bar)
	yield("DeviceSignalWiFi3BarLockIcon", icons.DeviceSignalWiFi3BarLock)
	yield("DeviceSignalWiFi4BarIcon", icons.DeviceSignalWiFi4Bar)
	yield("DeviceSignalWiFi4BarLockIcon", icons.DeviceSignalWiFi4BarLock)
	yield("DeviceSignalWiFiOffIcon", icons.DeviceSignalWiFiOff)
	yield("DeviceStorageIcon", icons.DeviceStorage)
	yield("DeviceUSBIcon", icons.DeviceUSB)
	yield("DeviceWallpaperIcon", icons.DeviceWallpaper)
	yield("DeviceWidgetsIcon", icons.DeviceWidgets)
	yield("DeviceWiFiLockIcon", icons.DeviceWiFiLock)
	yield("DeviceWiFiTetheringIcon", icons.DeviceWiFiTethering)
	yield("EditorAttachFileIcon", icons.EditorAttachFile)
	yield("EditorAttachMoneyIcon", icons.EditorAttachMoney)
	yield("EditorBorderAllIcon", icons.EditorBorderAll)
	yield("EditorBorderBottomIcon", icons.EditorBorderBottom)
	yield("EditorBorderClearIcon", icons.EditorBorderClear)
	yield("EditorBorderColorIcon", icons.EditorBorderColor)
	yield("EditorBorderHorizontalIcon", icons.EditorBorderHorizontal)
	yield("EditorBorderInnerIcon", icons.EditorBorderInner)
	yield("EditorBorderLeftIcon", icons.EditorBorderLeft)
	yield("EditorBorderOuterIcon", icons.EditorBorderOuter)
	yield("EditorBorderRightIcon", icons.EditorBorderRight)
	yield("EditorBorderStyleIcon", icons.EditorBorderStyle)
	yield("EditorBorderTopIcon", icons.EditorBorderTop)
	yield("EditorBorderVerticalIcon", icons.EditorBorderVertical)
	yield("EditorBubbleChartIcon", icons.EditorBubbleChart)
	yield("EditorDragHandleIcon", icons.EditorDragHandle)
	yield("EditorFormatAlignCenterIcon", icons.EditorFormatAlignCenter)
	yield("EditorFormatAlignJustifyIcon", icons.EditorFormatAlignJustify)
	yield("EditorFormatAlignLeftIcon", icons.EditorFormatAlignLeft)
	yield("EditorFormatAlignRightIcon", icons.EditorFormatAlignRight)
	yield("EditorFormatBoldIcon", icons.EditorFormatBold)
	yield("EditorFormatClearIcon", icons.EditorFormatClear)
	yield("EditorFormatColorFillIcon", icons.EditorFormatColorFill)
	yield("EditorFormatColorResetIcon", icons.EditorFormatColorReset)
	yield("EditorFormatColorTextIcon", icons.EditorFormatColorText)
	yield("EditorFormatIndentDecreaseIcon", icons.EditorFormatIndentDecrease)
	yield("EditorFormatIndentIncreaseIcon", icons.EditorFormatIndentIncrease)
	yield("EditorFormatItalicIcon", icons.EditorFormatItalic)
	yield("EditorFormatLineSpacingIcon", icons.EditorFormatLineSpacing)
	yield("EditorFormatListBulletedIcon", icons.EditorFormatListBulleted)
	yield("EditorFormatListNumberedIcon", icons.EditorFormatListNumbered)
	yield("EditorFormatPaintIcon", icons.EditorFormatPaint)
	yield("EditorFormatQuoteIcon", icons.EditorFormatQuote)
	yield("EditorFormatShapesIcon", icons.EditorFormatShapes)
	yield("EditorFormatSizeIcon", icons.EditorFormatSize)
	yield("EditorFormatStrikethroughIcon", icons.EditorFormatStrikethrough)
	yield("EditorFormatTextDirectionLToRIcon", icons.EditorFormatTextDirectionLToR)
	yield("EditorFormatTextDirectionRToLIcon", icons.EditorFormatTextDirectionRToL)
	yield("EditorFormatUnderlinedIcon", icons.EditorFormatUnderlined)
	yield("EditorFunctionsIcon", icons.EditorFunctions)
	yield("EditorHighlightIcon", icons.EditorHighlight)
	yield("EditorInsertChartIcon", icons.EditorInsertChart)
	yield("EditorInsertCommentIcon", icons.EditorInsertComment)
	yield("EditorInsertDriveFileIcon", icons.EditorInsertDriveFile)
	yield("EditorInsertEmoticonIcon", icons.EditorInsertEmoticon)
	yield("EditorInsertInvitationIcon", icons.EditorInsertInvitation)
	yield("EditorInsertLinkIcon", icons.EditorInsertLink)
	yield("EditorInsertPhotoIcon", icons.EditorInsertPhoto)
	yield("EditorLinearScaleIcon", icons.EditorLinearScale)
	yield("EditorMergeTypeIcon", icons.EditorMergeType)
	yield("EditorModeCommentIcon", icons.EditorModeComment)
	yield("EditorModeEditIcon", icons.EditorModeEdit)
	yield("EditorMonetizationOnIcon", icons.EditorMonetizationOn)
	yield("EditorMoneyOffIcon", icons.EditorMoneyOff)
	yield("EditorMultilineChartIcon", icons.EditorMultilineChart)
	yield("EditorPieChartIcon", icons.EditorPieChart)
	yield("EditorPieChartOutlinedIcon", icons.EditorPieChartOutlined)
	yield("EditorPublishIcon", icons.EditorPublish)
	yield("EditorShortTextIcon", icons.EditorShortText)
	yield("EditorShowChartIcon", icons.EditorShowChart)
	yield("EditorSpaceBarIcon", icons.EditorSpaceBar)
	yield("EditorStrikethroughSIcon", icons.EditorStrikethroughS)
	yield("EditorTextFieldsIcon", icons.EditorTextFields)
	yield("EditorTitleIcon", icons.EditorTitle)
	yield("EditorVerticalAlignBottomIcon", icons.EditorVerticalAlignBottom)
	yield("EditorVerticalAlignCenterIcon", icons.EditorVerticalAlignCenter)
	yield("EditorVerticalAlignTopIcon", icons.EditorVerticalAlignTop)
	yield("EditorWrapTextIcon", icons.EditorWrapText)
	yield("FileAttachmentIcon", icons.FileAttachment)
	yield("FileCloudIcon", icons.FileCloud)
	yield("FileCloudCircleIcon", icons.FileCloudCircle)
	yield("FileCloudDoneIcon", icons.FileCloudDone)
	yield("FileCloudDownloadIcon", icons.FileCloudDownload)
	yield("FileCloudOffIcon", icons.FileCloudOff)
	yield("FileCloudQueueIcon", icons.FileCloudQueue)
	yield("FileCloudUploadIcon", icons.FileCloudUpload)
	yield("FileCreateNewFolderIcon", icons.FileCreateNewFolder)
	yield("FileFileDownloadIcon", icons.FileFileDownload)
	yield("FileFileUploadIcon", icons.FileFileUpload)
	yield("FileFolderIcon", icons.FileFolder)
	yield("FileFolderOpenIcon", icons.FileFolderOpen)
	yield("FileFolderSharedIcon", icons.FileFolderShared)
	yield("HardwareCastIcon", icons.HardwareCast)
	yield("HardwareCastConnectedIcon", icons.HardwareCastConnected)
	yield("HardwareComputerIcon", icons.HardwareComputer)
	yield("HardwareDesktopMacIcon", icons.HardwareDesktopMac)
	yield("HardwareDesktopWindowsIcon", icons.HardwareDesktopWindows)
	yield("HardwareDeveloperBoardIcon", icons.HardwareDeveloperBoard)
	yield("HardwareDeviceHubIcon", icons.HardwareDeviceHub)
	yield("HardwareDevicesOtherIcon", icons.HardwareDevicesOther)
	yield("HardwareDockIcon", icons.HardwareDock)
	yield("HardwareGamepadIcon", icons.HardwareGamepad)
	yield("HardwareHeadsetIcon", icons.HardwareHeadset)
	yield("HardwareHeadsetMicIcon", icons.HardwareHeadsetMic)
	yield("HardwareKeyboardIcon", icons.HardwareKeyboard)
	yield("HardwareKeyboardArrowDownIcon", icons.HardwareKeyboardArrowDown)
	yield("HardwareKeyboardArrowLeftIcon", icons.HardwareKeyboardArrowLeft)
	yield("HardwareKeyboardArrowRightIcon", icons.HardwareKeyboardArrowRight)
	yield("HardwareKeyboardArrowUpIcon", icons.HardwareKeyboardArrowUp)
	yield("HardwareKeyboardBackspaceIcon", icons.HardwareKeyboardBackspace)
	yield("HardwareKeyboardCapslockIcon", icons.HardwareKeyboardCapslock)
	yield("HardwareKeyboardHideIcon", icons.HardwareKeyboardHide)
	yield("HardwareKeyboardReturnIcon", icons.HardwareKeyboardReturn)
	yield("HardwareKeyboardTabIcon", icons.HardwareKeyboardTab)
	yield("HardwareKeyboardVoiceIcon", icons.HardwareKeyboardVoice)
	yield("HardwareLaptopIcon", icons.HardwareLaptop)
	yield("HardwareLaptopChromebookIcon", icons.HardwareLaptopChromebook)
	yield("HardwareLaptopMacIcon", icons.HardwareLaptopMac)
	yield("HardwareLaptopWindowsIcon", icons.HardwareLaptopWindows)
	yield("HardwareMemoryIcon", icons.HardwareMemory)
	yield("HardwareMouseIcon", icons.HardwareMouse)
	yield("HardwarePhoneAndroidIcon", icons.HardwarePhoneAndroid)
	yield("HardwarePhoneIPhoneIcon", icons.HardwarePhoneIPhone)
	yield("HardwarePhoneLinkIcon", icons.HardwarePhoneLink)
	yield("HardwarePhoneLinkOffIcon", icons.HardwarePhoneLinkOff)
	yield("HardwarePowerInputIcon", icons.HardwarePowerInput)
	yield("HardwareRouterIcon", icons.HardwareRouter)
	yield("HardwareScannerIcon", icons.HardwareScanner)
	yield("HardwareSecurityIcon", icons.HardwareSecurity)
	yield("HardwareSIMCardIcon", icons.HardwareSIMCard)
	yield("HardwareSmartphoneIcon", icons.HardwareSmartphone)
	yield("HardwareSpeakerIcon", icons.HardwareSpeaker)
	yield("HardwareSpeakerGroupIcon", icons.HardwareSpeakerGroup)
	yield("HardwareTabletIcon", icons.HardwareTablet)
	yield("HardwareTabletAndroidIcon", icons.HardwareTabletAndroid)
	yield("HardwareTabletMacIcon", icons.HardwareTabletMac)
	yield("HardwareToysIcon", icons.HardwareToys)
	yield("HardwareTVIcon", icons.HardwareTV)
	yield("HardwareVideogameAssetIcon", icons.HardwareVideogameAsset)
	yield("HardwareWatchIcon", icons.HardwareWatch)
	yield("ImageAddAPhotoIcon", icons.ImageAddAPhoto)
	yield("ImageAddToPhotosIcon", icons.ImageAddToPhotos)
	yield("ImageAdjustIcon", icons.ImageAdjust)
	yield("ImageAssistantIcon", icons.ImageAssistant)
	yield("ImageAssistantPhotoIcon", icons.ImageAssistantPhoto)
	yield("ImageAudiotrackIcon", icons.ImageAudiotrack)
	yield("ImageBlurCircularIcon", icons.ImageBlurCircular)
	yield("ImageBlurLinearIcon", icons.ImageBlurLinear)
	yield("ImageBlurOffIcon", icons.ImageBlurOff)
	yield("ImageBlurOnIcon", icons.ImageBlurOn)
	yield("ImageBrightness1Icon", icons.ImageBrightness1)
	yield("ImageBrightness2Icon", icons.ImageBrightness2)
	yield("ImageBrightness3Icon", icons.ImageBrightness3)
	yield("ImageBrightness4Icon", icons.ImageBrightness4)
	yield("ImageBrightness5Icon", icons.ImageBrightness5)
	yield("ImageBrightness6Icon", icons.ImageBrightness6)
	yield("ImageBrightness7Icon", icons.ImageBrightness7)
	yield("ImageBrokenImageIcon", icons.ImageBrokenImage)
	yield("ImageBrushIcon", icons.ImageBrush)
	yield("ImageBurstModeIcon", icons.ImageBurstMode)
	yield("ImageCameraIcon", icons.ImageCamera)
	yield("ImageCameraAltIcon", icons.ImageCameraAlt)
	yield("ImageCameraFrontIcon", icons.ImageCameraFront)
	yield("ImageCameraRearIcon", icons.ImageCameraRear)
	yield("ImageCameraRollIcon", icons.ImageCameraRoll)
	yield("ImageCenterFocusStrongIcon", icons.ImageCenterFocusStrong)
	yield("ImageCenterFocusWeakIcon", icons.ImageCenterFocusWeak)
	yield("ImageCollectionsIcon", icons.ImageCollections)
	yield("ImageCollectionsBookmarkIcon", icons.ImageCollectionsBookmark)
	yield("ImageColorLensIcon", icons.ImageColorLens)
	yield("ImageColorizeIcon", icons.ImageColorize)
	yield("ImageCompareIcon", icons.ImageCompare)
	yield("ImageControlPointIcon", icons.ImageControlPoint)
	yield("ImageControlPointDuplicateIcon", icons.ImageControlPointDuplicate)
	yield("ImageCropIcon", icons.ImageCrop)
	yield("ImageCrop169Icon", icons.ImageCrop169)
	yield("ImageCrop32Icon", icons.ImageCrop32)
	yield("ImageCrop54Icon", icons.ImageCrop54)
	yield("ImageCrop75Icon", icons.ImageCrop75)
	yield("ImageCropDINIcon", icons.ImageCropDIN)
	yield("ImageCropFreeIcon", icons.ImageCropFree)
	yield("ImageCropLandscapeIcon", icons.ImageCropLandscape)
	yield("ImageCropOriginalIcon", icons.ImageCropOriginal)
	yield("ImageCropPortraitIcon", icons.ImageCropPortrait)
	yield("ImageCropRotateIcon", icons.ImageCropRotate)
	yield("ImageCropSquareIcon", icons.ImageCropSquare)
	yield("ImageDehazeIcon", icons.ImageDehaze)
	yield("ImageDetailsIcon", icons.ImageDetails)
	yield("ImageEditIcon", icons.ImageEdit)
	yield("ImageExposureIcon", icons.ImageExposure)
	yield("ImageExposureNeg1Icon", icons.ImageExposureNeg1)
	yield("ImageExposureNeg2Icon", icons.ImageExposureNeg2)
	yield("ImageExposurePlus1Icon", icons.ImageExposurePlus1)
	yield("ImageExposurePlus2Icon", icons.ImageExposurePlus2)
	yield("ImageExposureZeroIcon", icons.ImageExposureZero)
	yield("ImageFilterIcon", icons.ImageFilter)
	yield("ImageFilter1Icon", icons.ImageFilter1)
	yield("ImageFilter2Icon", icons.ImageFilter2)
	yield("ImageFilter3Icon", icons.ImageFilter3)
	yield("ImageFilter4Icon", icons.ImageFilter4)
	yield("ImageFilter5Icon", icons.ImageFilter5)
	yield("ImageFilter6Icon", icons.ImageFilter6)
	yield("ImageFilter7Icon", icons.ImageFilter7)
	yield("ImageFilter8Icon", icons.ImageFilter8)
	yield("ImageFilter9Icon", icons.ImageFilter9)
	yield("ImageFilter9PlusIcon", icons.ImageFilter9Plus)
	yield("ImageFilterBAndWIcon", icons.ImageFilterBAndW)
	yield("ImageFilterCenterFocusIcon", icons.ImageFilterCenterFocus)
	yield("ImageFilterDramaIcon", icons.ImageFilterDrama)
	yield("ImageFilterFramesIcon", icons.ImageFilterFrames)
	yield("ImageFilterHDRIcon", icons.ImageFilterHDR)
	yield("ImageFilterNoneIcon", icons.ImageFilterNone)
	yield("ImageFilterTiltShiftIcon", icons.ImageFilterTiltShift)
	yield("ImageFilterVintageIcon", icons.ImageFilterVintage)
	yield("ImageFlareIcon", icons.ImageFlare)
	yield("ImageFlashAutoIcon", icons.ImageFlashAuto)
	yield("ImageFlashOffIcon", icons.ImageFlashOff)
	yield("ImageFlashOnIcon", icons.ImageFlashOn)
	yield("ImageFlipIcon", icons.ImageFlip)
	yield("ImageGradientIcon", icons.ImageGradient)
	yield("ImageGrainIcon", icons.ImageGrain)
	yield("ImageGridOffIcon", icons.ImageGridOff)
	yield("ImageGridOnIcon", icons.ImageGridOn)
	yield("ImageHDROffIcon", icons.ImageHDROff)
	yield("ImageHDROnIcon", icons.ImageHDROn)
	yield("ImageHDRStrongIcon", icons.ImageHDRStrong)
	yield("ImageHDRWeakIcon", icons.ImageHDRWeak)
	yield("ImageHealingIcon", icons.ImageHealing)
	yield("ImageImageIcon", icons.ImageImage)
	yield("ImageImageAspectRatioIcon", icons.ImageImageAspectRatio)
	yield("ImageISOIcon", icons.ImageISO)
	yield("ImageLandscapeIcon", icons.ImageLandscape)
	yield("ImageLeakAddIcon", icons.ImageLeakAdd)
	yield("ImageLeakRemoveIcon", icons.ImageLeakRemove)
	yield("ImageLensIcon", icons.ImageLens)
	yield("ImageLinkedCameraIcon", icons.ImageLinkedCamera)
	yield("ImageLooksIcon", icons.ImageLooks)
	yield("ImageLooks3Icon", icons.ImageLooks3)
	yield("ImageLooks4Icon", icons.ImageLooks4)
	yield("ImageLooks5Icon", icons.ImageLooks5)
	yield("ImageLooks6Icon", icons.ImageLooks6)
	yield("ImageLooksOneIcon", icons.ImageLooksOne)
	yield("ImageLooksTwoIcon", icons.ImageLooksTwo)
	yield("ImageLoupeIcon", icons.ImageLoupe)
	yield("ImageMonochromePhotosIcon", icons.ImageMonochromePhotos)
	yield("ImageMovieCreationIcon", icons.ImageMovieCreation)
	yield("ImageMovieFilterIcon", icons.ImageMovieFilter)
	yield("ImageMusicNoteIcon", icons.ImageMusicNote)
	yield("ImageNatureIcon", icons.ImageNature)
	yield("ImageNaturePeopleIcon", icons.ImageNaturePeople)
	yield("ImageNavigateBeforeIcon", icons.ImageNavigateBefore)
	yield("ImageNavigateNextIcon", icons.ImageNavigateNext)
	yield("ImagePaletteIcon", icons.ImagePalette)
	yield("ImagePanoramaIcon", icons.ImagePanorama)
	yield("ImagePanoramaFishEyeIcon", icons.ImagePanoramaFishEye)
	yield("ImagePanoramaHorizontalIcon", icons.ImagePanoramaHorizontal)
	yield("ImagePanoramaVerticalIcon", icons.ImagePanoramaVertical)
	yield("ImagePanoramaWideAngleIcon", icons.ImagePanoramaWideAngle)
	yield("ImagePhotoIcon", icons.ImagePhoto)
	yield("ImagePhotoAlbumIcon", icons.ImagePhotoAlbum)
	yield("ImagePhotoCameraIcon", icons.ImagePhotoCamera)
	yield("ImagePhotoFilterIcon", icons.ImagePhotoFilter)
	yield("ImagePhotoLibraryIcon", icons.ImagePhotoLibrary)
	yield("ImagePhotoSizeSelectActualIcon", icons.ImagePhotoSizeSelectActual)
	yield("ImagePhotoSizeSelectLargeIcon", icons.ImagePhotoSizeSelectLarge)
	yield("ImagePhotoSizeSelectSmallIcon", icons.ImagePhotoSizeSelectSmall)
	yield("ImagePictureAsPDFIcon", icons.ImagePictureAsPDF)
	yield("ImagePortraitIcon", icons.ImagePortrait)
	yield("ImageRemoveRedEyeIcon", icons.ImageRemoveRedEye)
	yield("ImageRotate90DegreesCCWIcon", icons.ImageRotate90DegreesCCW)
	yield("ImageRotateLeftIcon", icons.ImageRotateLeft)
	yield("ImageRotateRightIcon", icons.ImageRotateRight)
	yield("ImageSlideshowIcon", icons.ImageSlideshow)
	yield("ImageStraightenIcon", icons.ImageStraighten)
	yield("ImageStyleIcon", icons.ImageStyle)
	yield("ImageSwitchCameraIcon", icons.ImageSwitchCamera)
	yield("ImageSwitchVideoIcon", icons.ImageSwitchVideo)
	yield("ImageTagFacesIcon", icons.ImageTagFaces)
	yield("ImageTextureIcon", icons.ImageTexture)
	yield("ImageTimeLapseIcon", icons.ImageTimeLapse)
	yield("ImageTimerIcon", icons.ImageTimer)
	yield("ImageTimer10Icon", icons.ImageTimer10)
	yield("ImageTimer3Icon", icons.ImageTimer3)
	yield("ImageTimerOffIcon", icons.ImageTimerOff)
	yield("ImageTonalityIcon", icons.ImageTonality)
	yield("ImageTransformIcon", icons.ImageTransform)
	yield("ImageTuneIcon", icons.ImageTune)
	yield("ImageViewComfyIcon", icons.ImageViewComfy)
	yield("ImageViewCompactIcon", icons.ImageViewCompact)
	yield("ImageVignetteIcon", icons.ImageVignette)
	yield("ImageWBAutoIcon", icons.ImageWBAuto)
	yield("ImageWBCloudyIcon", icons.ImageWBCloudy)
	yield("ImageWBIncandescentIcon", icons.ImageWBIncandescent)
	yield("ImageWBIridescentIcon", icons.ImageWBIridescent)
	yield("ImageWBSunnyIcon", icons.ImageWBSunny)
	yield("MapsAddLocationIcon", icons.MapsAddLocation)
	yield("MapsBeenhereIcon", icons.MapsBeenhere)
	yield("MapsDirectionsIcon", icons.MapsDirections)
	yield("MapsDirectionsBikeIcon", icons.MapsDirectionsBike)
	yield("MapsDirectionsBoatIcon", icons.MapsDirectionsBoat)
	yield("MapsDirectionsBusIcon", icons.MapsDirectionsBus)
	yield("MapsDirectionsCarIcon", icons.MapsDirectionsCar)
	yield("MapsDirectionsRailwayIcon", icons.MapsDirectionsRailway)
	yield("MapsDirectionsRunIcon", icons.MapsDirectionsRun)
	yield("MapsDirectionsSubwayIcon", icons.MapsDirectionsSubway)
	yield("MapsDirectionsTransitIcon", icons.MapsDirectionsTransit)
	yield("MapsDirectionsWalkIcon", icons.MapsDirectionsWalk)
	yield("MapsEditLocationIcon", icons.MapsEditLocation)
	yield("MapsEVStationIcon", icons.MapsEVStation)
	yield("MapsFlightIcon", icons.MapsFlight)
	yield("MapsHotelIcon", icons.MapsHotel)
	yield("MapsLayersIcon", icons.MapsLayers)
	yield("MapsLayersClearIcon", icons.MapsLayersClear)
	yield("MapsLocalActivityIcon", icons.MapsLocalActivity)
	yield("MapsLocalAirportIcon", icons.MapsLocalAirport)
	yield("MapsLocalATMIcon", icons.MapsLocalATM)
	yield("MapsLocalBarIcon", icons.MapsLocalBar)
	yield("MapsLocalCafeIcon", icons.MapsLocalCafe)
	yield("MapsLocalCarWashIcon", icons.MapsLocalCarWash)
	yield("MapsLocalConvenienceStoreIcon", icons.MapsLocalConvenienceStore)
	yield("MapsLocalDiningIcon", icons.MapsLocalDining)
	yield("MapsLocalDrinkIcon", icons.MapsLocalDrink)
	yield("MapsLocalFloristIcon", icons.MapsLocalFlorist)
	yield("MapsLocalGasStationIcon", icons.MapsLocalGasStation)
	yield("MapsLocalGroceryStoreIcon", icons.MapsLocalGroceryStore)
	yield("MapsLocalHospitalIcon", icons.MapsLocalHospital)
	yield("MapsLocalHotelIcon", icons.MapsLocalHotel)
	yield("MapsLocalLaundryServiceIcon", icons.MapsLocalLaundryService)
	yield("MapsLocalLibraryIcon", icons.MapsLocalLibrary)
	yield("MapsLocalMallIcon", icons.MapsLocalMall)
	yield("MapsLocalMoviesIcon", icons.MapsLocalMovies)
	yield("MapsLocalOfferIcon", icons.MapsLocalOffer)
	yield("MapsLocalParkingIcon", icons.MapsLocalParking)
	yield("MapsLocalPharmacyIcon", icons.MapsLocalPharmacy)
	yield("MapsLocalPhoneIcon", icons.MapsLocalPhone)
	yield("MapsLocalPizzaIcon", icons.MapsLocalPizza)
	yield("MapsLocalPlayIcon", icons.MapsLocalPlay)
	yield("MapsLocalPostOfficeIcon", icons.MapsLocalPostOffice)
	yield("MapsLocalPrintshopIcon", icons.MapsLocalPrintshop)
	yield("MapsLocalSeeIcon", icons.MapsLocalSee)
	yield("MapsLocalShippingIcon", icons.MapsLocalShipping)
	yield("MapsLocalTaxiIcon", icons.MapsLocalTaxi)
	yield("MapsMapIcon", icons.MapsMap)
	yield("MapsMyLocationIcon", icons.MapsMyLocation)
	yield("MapsNavigationIcon", icons.MapsNavigation)
	yield("MapsNearMeIcon", icons.MapsNearMe)
	yield("MapsPersonPinIcon", icons.MapsPersonPin)
	yield("MapsPersonPinCircleIcon", icons.MapsPersonPinCircle)
	yield("MapsPinDropIcon", icons.MapsPinDrop)
	yield("MapsPlaceIcon", icons.MapsPlace)
	yield("MapsRateReviewIcon", icons.MapsRateReview)
	yield("MapsRestaurantIcon", icons.MapsRestaurant)
	yield("MapsRestaurantMenuIcon", icons.MapsRestaurantMenu)
	yield("MapsSatelliteIcon", icons.MapsSatellite)
	yield("MapsStoreMallDirectoryIcon", icons.MapsStoreMallDirectory)
	yield("MapsStreetViewIcon", icons.MapsStreetView)
	yield("MapsSubwayIcon", icons.MapsSubway)
	yield("MapsTerrainIcon", icons.MapsTerrain)
	yield("MapsTrafficIcon", icons.MapsTraffic)
	yield("MapsTrainIcon", icons.MapsTrain)
	yield("MapsTramIcon", icons.MapsTram)
	yield("MapsTransferWithinAStationIcon", icons.MapsTransferWithinAStation)
	yield("MapsZoomOutMapIcon", icons.MapsZoomOutMap)
	yield("NavigationAppsIcon", icons.NavigationApps)
	yield("NavigationArrowBackIcon", icons.NavigationArrowBack)
	yield("NavigationArrowDownwardIcon", icons.NavigationArrowDownward)
	yield("NavigationArrowDropDownIcon", icons.NavigationArrowDropDown)
	yield("NavigationArrowDropDownCircleIcon", icons.NavigationArrowDropDownCircle)
	yield("NavigationArrowDropUpIcon", icons.NavigationArrowDropUp)
	yield("NavigationArrowForwardIcon", icons.NavigationArrowForward)
	yield("NavigationArrowUpwardIcon", icons.NavigationArrowUpward)
	yield("NavigationCancelIcon", icons.NavigationCancel)
	yield("NavigationCheckIcon", icons.NavigationCheck)
	yield("NavigationChevronLeftIcon", icons.NavigationChevronLeft)
	yield("NavigationChevronRightIcon", icons.NavigationChevronRight)
	yield("NavigationCloseIcon", icons.NavigationClose)
	yield("NavigationExpandLessIcon", icons.NavigationExpandLess)
	yield("NavigationExpandMoreIcon", icons.NavigationExpandMore)
	yield("NavigationFirstPageIcon", icons.NavigationFirstPage)
	yield("NavigationFullscreenIcon", icons.NavigationFullscreen)
	yield("NavigationFullscreenExitIcon", icons.NavigationFullscreenExit)
	yield("NavigationLastPageIcon", icons.NavigationLastPage)
	yield("NavigationMenuIcon", icons.NavigationMenu)
	yield("NavigationMoreHorizIcon", icons.NavigationMoreHoriz)
	yield("NavigationMoreVertIcon", icons.NavigationMoreVert)
	yield("NavigationRefreshIcon", icons.NavigationRefresh)
	yield("NavigationSubdirectoryArrowLeftIcon", icons.NavigationSubdirectoryArrowLeft)
	yield("NavigationSubdirectoryArrowRightIcon", icons.NavigationSubdirectoryArrowRight)
	yield("NavigationUnfoldLessIcon", icons.NavigationUnfoldLess)
	yield("NavigationUnfoldMoreIcon", icons.NavigationUnfoldMore)
	yield("NotificationADBIcon", icons.NotificationADB)
	yield("NotificationAirlineSeatFlatIcon", icons.NotificationAirlineSeatFlat)
	yield("NotificationAirlineSeatFlatAngledIcon", icons.NotificationAirlineSeatFlatAngled)
	yield("NotificationAirlineSeatIndividualSuiteIcon", icons.NotificationAirlineSeatIndividualSuite)
	yield("NotificationAirlineSeatLegroomExtraIcon", icons.NotificationAirlineSeatLegroomExtra)
	yield("NotificationAirlineSeatLegroomNormalIcon", icons.NotificationAirlineSeatLegroomNormal)
	yield("NotificationAirlineSeatLegroomReducedIcon", icons.NotificationAirlineSeatLegroomReduced)
	yield("NotificationAirlineSeatReclineExtraIcon", icons.NotificationAirlineSeatReclineExtra)
	yield("NotificationAirlineSeatReclineNormalIcon", icons.NotificationAirlineSeatReclineNormal)
	yield("NotificationBluetoothAudioIcon", icons.NotificationBluetoothAudio)
	yield("NotificationConfirmationNumberIcon", icons.NotificationConfirmationNumber)
	yield("NotificationDiscFullIcon", icons.NotificationDiscFull)
	yield("NotificationDoNotDisturbIcon", icons.NotificationDoNotDisturb)
	yield("NotificationDoNotDisturbAltIcon", icons.NotificationDoNotDisturbAlt)
	yield("NotificationDoNotDisturbOffIcon", icons.NotificationDoNotDisturbOff)
	yield("NotificationDoNotDisturbOnIcon", icons.NotificationDoNotDisturbOn)
	yield("NotificationDriveETAIcon", icons.NotificationDriveETA)
	yield("NotificationEnhancedEncryptionIcon", icons.NotificationEnhancedEncryption)
	yield("NotificationEventAvailableIcon", icons.NotificationEventAvailable)
	yield("NotificationEventBusyIcon", icons.NotificationEventBusy)
	yield("NotificationEventNoteIcon", icons.NotificationEventNote)
	yield("NotificationFolderSpecialIcon", icons.NotificationFolderSpecial)
	yield("NotificationLiveTVIcon", icons.NotificationLiveTV)
	yield("NotificationMMSIcon", icons.NotificationMMS)
	yield("NotificationMoreIcon", icons.NotificationMore)
	yield("NotificationNetworkCheckIcon", icons.NotificationNetworkCheck)
	yield("NotificationNetworkLockedIcon", icons.NotificationNetworkLocked)
	yield("NotificationNoEncryptionIcon", icons.NotificationNoEncryption)
	yield("NotificationOnDemandVideoIcon", icons.NotificationOnDemandVideo)
	yield("NotificationPersonalVideoIcon", icons.NotificationPersonalVideo)
	yield("NotificationPhoneBluetoothSpeakerIcon", icons.NotificationPhoneBluetoothSpeaker)
	yield("NotificationPhoneForwardedIcon", icons.NotificationPhoneForwarded)
	yield("NotificationPhoneInTalkIcon", icons.NotificationPhoneInTalk)
	yield("NotificationPhoneLockedIcon", icons.NotificationPhoneLocked)
	yield("NotificationPhoneMissedIcon", icons.NotificationPhoneMissed)
	yield("NotificationPhonePausedIcon", icons.NotificationPhonePaused)
	yield("NotificationPowerIcon", icons.NotificationPower)
	yield("NotificationPriorityHighIcon", icons.NotificationPriorityHigh)
	yield("NotificationRVHookupIcon", icons.NotificationRVHookup)
	yield("NotificationSDCardIcon", icons.NotificationSDCard)
	yield("NotificationSIMCardAlertIcon", icons.NotificationSIMCardAlert)
	yield("NotificationSMSIcon", icons.NotificationSMS)
	yield("NotificationSMSFailedIcon", icons.NotificationSMSFailed)
	yield("NotificationSyncIcon", icons.NotificationSync)
	yield("NotificationSyncDisabledIcon", icons.NotificationSyncDisabled)
	yield("NotificationSyncProblemIcon", icons.NotificationSyncProblem)
	yield("NotificationSystemUpdateIcon", icons.NotificationSystemUpdate)
	yield("NotificationTapAndPlayIcon", icons.NotificationTapAndPlay)
	yield("NotificationTimeToLeaveIcon", icons.NotificationTimeToLeave)
	yield("NotificationVibrationIcon", icons.NotificationVibration)
	yield("NotificationVoiceChatIcon", icons.NotificationVoiceChat)
	yield("NotificationVPNLockIcon", icons.NotificationVPNLock)
	yield("NotificationWCIcon", icons.NotificationWC)
	yield("NotificationWiFiIcon", icons.NotificationWiFi)
	yield("PlacesACUnitIcon", icons.PlacesACUnit)
	yield("PlacesAirportShuttleIcon", icons.PlacesAirportShuttle)
	yield("PlacesAllInclusiveIcon", icons.PlacesAllInclusive)
	yield("PlacesBeachAccessIcon", icons.PlacesBeachAccess)
	yield("PlacesBusinessCenterIcon", icons.PlacesBusinessCenter)
	yield("PlacesCasinoIcon", icons.PlacesCasino)
	yield("PlacesChildCareIcon", icons.PlacesChildCare)
	yield("PlacesChildFriendlyIcon", icons.PlacesChildFriendly)
	yield("PlacesFitnessCenterIcon", icons.PlacesFitnessCenter)
	yield("PlacesFreeBreakfastIcon", icons.PlacesFreeBreakfast)
	yield("PlacesGolfCourseIcon", icons.PlacesGolfCourse)
	yield("PlacesHotTubIcon", icons.PlacesHotTub)
	yield("PlacesKitchenIcon", icons.PlacesKitchen)
	yield("PlacesPoolIcon", icons.PlacesPool)
	yield("PlacesRoomServiceIcon", icons.PlacesRoomService)
	yield("PlacesRVHookupIcon", icons.PlacesRVHookup)
	yield("PlacesSmokeFreeIcon", icons.PlacesSmokeFree)
	yield("PlacesSmokingRoomsIcon", icons.PlacesSmokingRooms)
	yield("PlacesSpaIcon", icons.PlacesSpa)
	yield("SocialCakeIcon", icons.SocialCake)
	yield("SocialDomainIcon", icons.SocialDomain)
	yield("SocialGroupIcon", icons.SocialGroup)
	yield("SocialGroupAddIcon", icons.SocialGroupAdd)
	yield("SocialLocationCityIcon", icons.SocialLocationCity)
	yield("SocialMoodIcon", icons.SocialMood)
	yield("SocialMoodBadIcon", icons.SocialMoodBad)
	yield("SocialNotificationsIcon", icons.SocialNotifications)
	yield("SocialNotificationsActiveIcon", icons.SocialNotificationsActive)
	yield("SocialNotificationsNoneIcon", icons.SocialNotificationsNone)
	yield("SocialNotificationsOffIcon", icons.SocialNotificationsOff)
	yield("SocialNotificationsPausedIcon", icons.SocialNotificationsPaused)
	yield("SocialPagesIcon", icons.SocialPages)
	yield("SocialPartyModeIcon", icons.SocialPartyMode)
	yield("SocialPeopleIcon", icons.SocialPeople)
	yield("SocialPeopleOutlineIcon", icons.SocialPeopleOutline)
	yield("SocialPersonIcon", icons.SocialPerson)
	yield("SocialPersonAddIcon", icons.SocialPersonAdd)
	yield("SocialPersonOutlineIcon", icons.SocialPersonOutline)
	yield("SocialPlusOneIcon", icons.SocialPlusOne)
	yield("SocialPollIcon", icons.SocialPoll)
	yield("SocialPublicIcon", icons.SocialPublic)
	yield("SocialSchoolIcon", icons.SocialSchool)
	yield("SocialSentimentDissatisfiedIcon", icons.SocialSentimentDissatisfied)
	yield("SocialSentimentNeutralIcon", icons.SocialSentimentNeutral)
	yield("SocialSentimentSatisfiedIcon", icons.SocialSentimentSatisfied)
	yield("SocialSentimentVeryDissatisfiedIcon", icons.SocialSentimentVeryDissatisfied)
	yield("SocialSentimentVerySatisfiedIcon", icons.SocialSentimentVerySatisfied)
	yield("SocialShareIcon", icons.SocialShare)
	yield("SocialWhatsHotIcon", icons.SocialWhatsHot)
	yield("ToggleCheckBoxIcon", icons.ToggleCheckBox)
	yield("ToggleCheckBoxOutlineBlankIcon", icons.ToggleCheckBoxOutlineBlank)
	yield("ToggleIndeterminateCheckBoxIcon", icons.ToggleIndeterminateCheckBox)
	yield("ToggleRadioButtonCheckedIcon", icons.ToggleRadioButtonChecked)
	yield("ToggleRadioButtonUncheckedIcon", icons.ToggleRadioButtonUnchecked)
	yield("ToggleStarIcon", icons.ToggleStar)
	yield("ToggleStarBorderIcon", icons.ToggleStarBorder)
	yield("ToggleStarHalfIcon", icons.ToggleStarHalf)
})

var (
	Action3DRotationIcon                            = icons.Action3DRotation
	ActionAccessibilityIcon                         = icons.ActionAccessibility
	ActionAccessibleIcon                            = icons.ActionAccessible
	ActionAccountBalanceIcon                        = icons.ActionAccountBalance
	ActionAccountBalanceWalletIcon                  = icons.ActionAccountBalanceWallet
	ActionAccountBoxIcon                            = icons.ActionAccountBox
	ActionAccountCircleIcon                         = icons.ActionAccountCircle
	ActionAddShoppingCartIcon                       = icons.ActionAddShoppingCart
	ActionAlarmIcon                                 = icons.ActionAlarm
	ActionAlarmAddIcon                              = icons.ActionAlarmAdd
	ActionAlarmOffIcon                              = icons.ActionAlarmOff
	ActionAlarmOnIcon                               = icons.ActionAlarmOn
	ActionAllOutIcon                                = icons.ActionAllOut
	ActionAndroidIcon                               = icons.ActionAndroid
	ActionAnnouncementIcon                          = icons.ActionAnnouncement
	ActionAspectRatioIcon                           = icons.ActionAspectRatio
	ActionAssessmentIcon                            = icons.ActionAssessment
	ActionAssignmentIcon                            = icons.ActionAssignment
	ActionAssignmentIndIcon                         = icons.ActionAssignmentInd
	ActionAssignmentLateIcon                        = icons.ActionAssignmentLate
	ActionAssignmentReturnIcon                      = icons.ActionAssignmentReturn
	ActionAssignmentReturnedIcon                    = icons.ActionAssignmentReturned
	ActionAssignmentTurnedInIcon                    = icons.ActionAssignmentTurnedIn
	ActionAutorenewIcon                             = icons.ActionAutorenew
	ActionBackupIcon                                = icons.ActionBackup
	ActionBookIcon                                  = icons.ActionBook
	ActionBookmarkIcon                              = icons.ActionBookmark
	ActionBookmarkBorderIcon                        = icons.ActionBookmarkBorder
	ActionBugReportIcon                             = icons.ActionBugReport
	ActionBuildIcon                                 = icons.ActionBuild
	ActionCachedIcon                                = icons.ActionCached
	ActionCameraEnhanceIcon                         = icons.ActionCameraEnhance
	ActionCardGiftcardIcon                          = icons.ActionCardGiftcard
	ActionCardMembershipIcon                        = icons.ActionCardMembership
	ActionCardTravelIcon                            = icons.ActionCardTravel
	ActionChangeHistoryIcon                         = icons.ActionChangeHistory
	ActionCheckCircleIcon                           = icons.ActionCheckCircle
	ActionChromeReaderModeIcon                      = icons.ActionChromeReaderMode
	ActionClassIcon                                 = icons.ActionClass
	ActionCodeIcon                                  = icons.ActionCode
	ActionCompareArrowsIcon                         = icons.ActionCompareArrows
	ActionCopyrightIcon                             = icons.ActionCopyright
	ActionCreditCardIcon                            = icons.ActionCreditCard
	ActionDashboardIcon                             = icons.ActionDashboard
	ActionDateRangeIcon                             = icons.ActionDateRange
	ActionDeleteIcon                                = icons.ActionDelete
	ActionDeleteForeverIcon                         = icons.ActionDeleteForever
	ActionDescriptionIcon                           = icons.ActionDescription
	ActionDNSIcon                                   = icons.ActionDNS
	ActionDoneIcon                                  = icons.ActionDone
	ActionDoneAllIcon                               = icons.ActionDoneAll
	ActionDonutLargeIcon                            = icons.ActionDonutLarge
	ActionDonutSmallIcon                            = icons.ActionDonutSmall
	ActionEjectIcon                                 = icons.ActionEject
	ActionEuroSymbolIcon                            = icons.ActionEuroSymbol
	ActionEventIcon                                 = icons.ActionEvent
	ActionEventSeatIcon                             = icons.ActionEventSeat
	ActionExitToAppIcon                             = icons.ActionExitToApp
	ActionExploreIcon                               = icons.ActionExplore
	ActionExtensionIcon                             = icons.ActionExtension
	ActionFaceIcon                                  = icons.ActionFace
	ActionFavoriteIcon                              = icons.ActionFavorite
	ActionFavoriteBorderIcon                        = icons.ActionFavoriteBorder
	ActionFeedbackIcon                              = icons.ActionFeedback
	ActionFindInPageIcon                            = icons.ActionFindInPage
	ActionFindReplaceIcon                           = icons.ActionFindReplace
	ActionFingerprintIcon                           = icons.ActionFingerprint
	ActionFlightLandIcon                            = icons.ActionFlightLand
	ActionFlightTakeoffIcon                         = icons.ActionFlightTakeoff
	ActionFlipToBackIcon                            = icons.ActionFlipToBack
	ActionFlipToFrontIcon                           = icons.ActionFlipToFront
	ActionGTranslateIcon                            = icons.ActionGTranslate
	ActionGavelIcon                                 = icons.ActionGavel
	ActionGetAppIcon                                = icons.ActionGetApp
	ActionGIFIcon                                   = icons.ActionGIF
	ActionGradeIcon                                 = icons.ActionGrade
	ActionGroupWorkIcon                             = icons.ActionGroupWork
	ActionHelpIcon                                  = icons.ActionHelp
	ActionHelpOutlineIcon                           = icons.ActionHelpOutline
	ActionHighlightOffIcon                          = icons.ActionHighlightOff
	ActionHistoryIcon                               = icons.ActionHistory
	ActionHomeIcon                                  = icons.ActionHome
	ActionHourglassEmptyIcon                        = icons.ActionHourglassEmpty
	ActionHourglassFullIcon                         = icons.ActionHourglassFull
	ActionHTTPIcon                                  = icons.ActionHTTP
	ActionHTTPSIcon                                 = icons.ActionHTTPS
	ActionImportantDevicesIcon                      = icons.ActionImportantDevices
	ActionInfoIcon                                  = icons.ActionInfo
	ActionInfoOutlineIcon                           = icons.ActionInfoOutline
	ActionInputIcon                                 = icons.ActionInput
	ActionInvertColorsIcon                          = icons.ActionInvertColors
	ActionLabelIcon                                 = icons.ActionLabel
	ActionLabelOutlineIcon                          = icons.ActionLabelOutline
	ActionLanguageIcon                              = icons.ActionLanguage
	ActionLaunchIcon                                = icons.ActionLaunch
	ActionLightbulbOutlineIcon                      = icons.ActionLightbulbOutline
	ActionLineStyleIcon                             = icons.ActionLineStyle
	ActionLineWeightIcon                            = icons.ActionLineWeight
	ActionListIcon                                  = icons.ActionList
	ActionLockIcon                                  = icons.ActionLock
	ActionLockOpenIcon                              = icons.ActionLockOpen
	ActionLockOutlineIcon                           = icons.ActionLockOutline
	ActionLoyaltyIcon                               = icons.ActionLoyalty
	ActionMarkUnreadMailboxIcon                     = icons.ActionMarkUnreadMailbox
	ActionMotorcycleIcon                            = icons.ActionMotorcycle
	ActionNoteAddIcon                               = icons.ActionNoteAdd
	ActionOfflinePinIcon                            = icons.ActionOfflinePin
	ActionOpacityIcon                               = icons.ActionOpacity
	ActionOpenInBrowserIcon                         = icons.ActionOpenInBrowser
	ActionOpenInNewIcon                             = icons.ActionOpenInNew
	ActionOpenWithIcon                              = icons.ActionOpenWith
	ActionPageviewIcon                              = icons.ActionPageview
	ActionPanToolIcon                               = icons.ActionPanTool
	ActionPaymentIcon                               = icons.ActionPayment
	ActionPermCameraMicIcon                         = icons.ActionPermCameraMic
	ActionPermContactCalendarIcon                   = icons.ActionPermContactCalendar
	ActionPermDataSettingIcon                       = icons.ActionPermDataSetting
	ActionPermDeviceInformationIcon                 = icons.ActionPermDeviceInformation
	ActionPermIdentityIcon                          = icons.ActionPermIdentity
	ActionPermMediaIcon                             = icons.ActionPermMedia
	ActionPermPhoneMsgIcon                          = icons.ActionPermPhoneMsg
	ActionPermScanWiFiIcon                          = icons.ActionPermScanWiFi
	ActionPetsIcon                                  = icons.ActionPets
	ActionPictureInPictureIcon                      = icons.ActionPictureInPicture
	ActionPictureInPictureAltIcon                   = icons.ActionPictureInPictureAlt
	ActionPlayForWorkIcon                           = icons.ActionPlayForWork
	ActionPolymerIcon                               = icons.ActionPolymer
	ActionPowerSettingsNewIcon                      = icons.ActionPowerSettingsNew
	ActionPregnantWomanIcon                         = icons.ActionPregnantWoman
	ActionPrintIcon                                 = icons.ActionPrint
	ActionQueryBuilderIcon                          = icons.ActionQueryBuilder
	ActionQuestionAnswerIcon                        = icons.ActionQuestionAnswer
	ActionReceiptIcon                               = icons.ActionReceipt
	ActionRecordVoiceOverIcon                       = icons.ActionRecordVoiceOver
	ActionRedeemIcon                                = icons.ActionRedeem
	ActionRemoveShoppingCartIcon                    = icons.ActionRemoveShoppingCart
	ActionReorderIcon                               = icons.ActionReorder
	ActionReportProblemIcon                         = icons.ActionReportProblem
	ActionRestoreIcon                               = icons.ActionRestore
	ActionRestorePageIcon                           = icons.ActionRestorePage
	ActionRoomIcon                                  = icons.ActionRoom
	ActionRoundedCornerIcon                         = icons.ActionRoundedCorner
	ActionRowingIcon                                = icons.ActionRowing
	ActionScheduleIcon                              = icons.ActionSchedule
	ActionSearchIcon                                = icons.ActionSearch
	ActionSettingsIcon                              = icons.ActionSettings
	ActionSettingsApplicationsIcon                  = icons.ActionSettingsApplications
	ActionSettingsBackupRestoreIcon                 = icons.ActionSettingsBackupRestore
	ActionSettingsBluetoothIcon                     = icons.ActionSettingsBluetooth
	ActionSettingsBrightnessIcon                    = icons.ActionSettingsBrightness
	ActionSettingsCellIcon                          = icons.ActionSettingsCell
	ActionSettingsEthernetIcon                      = icons.ActionSettingsEthernet
	ActionSettingsInputAntennaIcon                  = icons.ActionSettingsInputAntenna
	ActionSettingsInputComponentIcon                = icons.ActionSettingsInputComponent
	ActionSettingsInputCompositeIcon                = icons.ActionSettingsInputComposite
	ActionSettingsInputHDMIIcon                     = icons.ActionSettingsInputHDMI
	ActionSettingsInputSVideoIcon                   = icons.ActionSettingsInputSVideo
	ActionSettingsOverscanIcon                      = icons.ActionSettingsOverscan
	ActionSettingsPhoneIcon                         = icons.ActionSettingsPhone
	ActionSettingsPowerIcon                         = icons.ActionSettingsPower
	ActionSettingsRemoteIcon                        = icons.ActionSettingsRemote
	ActionSettingsVoiceIcon                         = icons.ActionSettingsVoice
	ActionShopIcon                                  = icons.ActionShop
	ActionShopTwoIcon                               = icons.ActionShopTwo
	ActionShoppingBasketIcon                        = icons.ActionShoppingBasket
	ActionShoppingCartIcon                          = icons.ActionShoppingCart
	ActionSpeakerNotesIcon                          = icons.ActionSpeakerNotes
	ActionSpeakerNotesOffIcon                       = icons.ActionSpeakerNotesOff
	ActionSpellcheckIcon                            = icons.ActionSpellcheck
	ActionStarRateIcon                              = icons.ActionStarRate
	ActionStarsIcon                                 = icons.ActionStars
	ActionStoreIcon                                 = icons.ActionStore
	ActionSubjectIcon                               = icons.ActionSubject
	ActionSupervisorAccountIcon                     = icons.ActionSupervisorAccount
	ActionSwapHorizIcon                             = icons.ActionSwapHoriz
	ActionSwapVertIcon                              = icons.ActionSwapVert
	ActionSwapVerticalCircleIcon                    = icons.ActionSwapVerticalCircle
	ActionSystemUpdateAltIcon                       = icons.ActionSystemUpdateAlt
	ActionTabIcon                                   = icons.ActionTab
	ActionTabUnselectedIcon                         = icons.ActionTabUnselected
	ActionTheatersIcon                              = icons.ActionTheaters
	ActionThumbDownIcon                             = icons.ActionThumbDown
	ActionThumbUpIcon                               = icons.ActionThumbUp
	ActionThumbsUpDownIcon                          = icons.ActionThumbsUpDown
	ActionTimelineIcon                              = icons.ActionTimeline
	ActionTOCIcon                                   = icons.ActionTOC
	ActionTodayIcon                                 = icons.ActionToday
	ActionTollIcon                                  = icons.ActionToll
	ActionTouchAppIcon                              = icons.ActionTouchApp
	ActionTrackChangesIcon                          = icons.ActionTrackChanges
	ActionTranslateIcon                             = icons.ActionTranslate
	ActionTrendingDownIcon                          = icons.ActionTrendingDown
	ActionTrendingFlatIcon                          = icons.ActionTrendingFlat
	ActionTrendingUpIcon                            = icons.ActionTrendingUp
	ActionTurnedInIcon                              = icons.ActionTurnedIn
	ActionTurnedInNotIcon                           = icons.ActionTurnedInNot
	ActionUpdateIcon                                = icons.ActionUpdate
	ActionVerifiedUserIcon                          = icons.ActionVerifiedUser
	ActionViewAgendaIcon                            = icons.ActionViewAgenda
	ActionViewArrayIcon                             = icons.ActionViewArray
	ActionViewCarouselIcon                          = icons.ActionViewCarousel
	ActionViewColumnIcon                            = icons.ActionViewColumn
	ActionViewDayIcon                               = icons.ActionViewDay
	ActionViewHeadlineIcon                          = icons.ActionViewHeadline
	ActionViewListIcon                              = icons.ActionViewList
	ActionViewModuleIcon                            = icons.ActionViewModule
	ActionViewQuiltIcon                             = icons.ActionViewQuilt
	ActionViewStreamIcon                            = icons.ActionViewStream
	ActionViewWeekIcon                              = icons.ActionViewWeek
	ActionVisibilityIcon                            = icons.ActionVisibility
	ActionVisibilityOffIcon                         = icons.ActionVisibilityOff
	ActionWatchLaterIcon                            = icons.ActionWatchLater
	ActionWorkIcon                                  = icons.ActionWork
	ActionYoutubeSearchedForIcon                    = icons.ActionYoutubeSearchedFor
	ActionZoomInIcon                                = icons.ActionZoomIn
	ActionZoomOutIcon                               = icons.ActionZoomOut
	AlertAddAlertIcon                               = icons.AlertAddAlert
	AlertErrorIcon                                  = icons.AlertError
	AlertErrorOutlineIcon                           = icons.AlertErrorOutline
	AlertWarningIcon                                = icons.AlertWarning
	AVAddToQueueIcon                                = icons.AVAddToQueue
	AVAirplayIcon                                   = icons.AVAirplay
	AVAlbumIcon                                     = icons.AVAlbum
	AVArtTrackIcon                                  = icons.AVArtTrack
	AVAVTimerIcon                                   = icons.AVAVTimer
	AVBrandingWatermarkIcon                         = icons.AVBrandingWatermark
	AVCallToActionIcon                              = icons.AVCallToAction
	AVClosedCaptionIcon                             = icons.AVClosedCaption
	AVEqualizerIcon                                 = icons.AVEqualizer
	AVExplicitIcon                                  = icons.AVExplicit
	AVFastForwardIcon                               = icons.AVFastForward
	AVFastRewindIcon                                = icons.AVFastRewind
	AVFeaturedPlayListIcon                          = icons.AVFeaturedPlayList
	AVFeaturedVideoIcon                             = icons.AVFeaturedVideo
	AVFiberDVRIcon                                  = icons.AVFiberDVR
	AVFiberManualRecordIcon                         = icons.AVFiberManualRecord
	AVFiberNewIcon                                  = icons.AVFiberNew
	AVFiberPinIcon                                  = icons.AVFiberPin
	AVFiberSmartRecordIcon                          = icons.AVFiberSmartRecord
	AVForward10Icon                                 = icons.AVForward10
	AVForward30Icon                                 = icons.AVForward30
	AVForward5Icon                                  = icons.AVForward5
	AVGamesIcon                                     = icons.AVGames
	AVHDIcon                                        = icons.AVHD
	AVHearingIcon                                   = icons.AVHearing
	AVHighQualityIcon                               = icons.AVHighQuality
	AVLibraryAddIcon                                = icons.AVLibraryAdd
	AVLibraryBooksIcon                              = icons.AVLibraryBooks
	AVLibraryMusicIcon                              = icons.AVLibraryMusic
	AVLoopIcon                                      = icons.AVLoop
	AVMicIcon                                       = icons.AVMic
	AVMicNoneIcon                                   = icons.AVMicNone
	AVMicOffIcon                                    = icons.AVMicOff
	AVMovieIcon                                     = icons.AVMovie
	AVMusicVideoIcon                                = icons.AVMusicVideo
	AVNewReleasesIcon                               = icons.AVNewReleases
	AVNotInterestedIcon                             = icons.AVNotInterested
	AVNoteIcon                                      = icons.AVNote
	AVPauseIcon                                     = icons.AVPause
	AVPauseCircleFilledIcon                         = icons.AVPauseCircleFilled
	AVPauseCircleOutlineIcon                        = icons.AVPauseCircleOutline
	AVPlayArrowIcon                                 = icons.AVPlayArrow
	AVPlayCircleFilledIcon                          = icons.AVPlayCircleFilled
	AVPlayCircleOutlineIcon                         = icons.AVPlayCircleOutline
	AVPlaylistAddIcon                               = icons.AVPlaylistAdd
	AVPlaylistAddCheckIcon                          = icons.AVPlaylistAddCheck
	AVPlaylistPlayIcon                              = icons.AVPlaylistPlay
	AVQueueIcon                                     = icons.AVQueue
	AVQueueMusicIcon                                = icons.AVQueueMusic
	AVQueuePlayNextIcon                             = icons.AVQueuePlayNext
	AVRadioIcon                                     = icons.AVRadio
	AVRecentActorsIcon                              = icons.AVRecentActors
	AVRemoveFromQueueIcon                           = icons.AVRemoveFromQueue
	AVRepeatIcon                                    = icons.AVRepeat
	AVRepeatOneIcon                                 = icons.AVRepeatOne
	AVReplayIcon                                    = icons.AVReplay
	AVReplay10Icon                                  = icons.AVReplay10
	AVReplay30Icon                                  = icons.AVReplay30
	AVReplay5Icon                                   = icons.AVReplay5
	AVShuffleIcon                                   = icons.AVShuffle
	AVSkipNextIcon                                  = icons.AVSkipNext
	AVSkipPreviousIcon                              = icons.AVSkipPrevious
	AVSlowMotionVideoIcon                           = icons.AVSlowMotionVideo
	AVSnoozeIcon                                    = icons.AVSnooze
	AVSortByAlphaIcon                               = icons.AVSortByAlpha
	AVStopIcon                                      = icons.AVStop
	AVSubscriptionsIcon                             = icons.AVSubscriptions
	AVSubtitlesIcon                                 = icons.AVSubtitles
	AVSurroundSoundIcon                             = icons.AVSurroundSound
	AVVideoCallIcon                                 = icons.AVVideoCall
	AVVideoLabelIcon                                = icons.AVVideoLabel
	AVVideoLibraryIcon                              = icons.AVVideoLibrary
	AVVideocamIcon                                  = icons.AVVideocam
	AVVideocamOffIcon                               = icons.AVVideocamOff
	AVVolumeDownIcon                                = icons.AVVolumeDown
	AVVolumeMuteIcon                                = icons.AVVolumeMute
	AVVolumeOffIcon                                 = icons.AVVolumeOff
	AVVolumeUpIcon                                  = icons.AVVolumeUp
	AVWebIcon                                       = icons.AVWeb
	AVWebAssetIcon                                  = icons.AVWebAsset
	CommunicationBusinessIcon                       = icons.CommunicationBusiness
	CommunicationCallIcon                           = icons.CommunicationCall
	CommunicationCallEndIcon                        = icons.CommunicationCallEnd
	CommunicationCallMadeIcon                       = icons.CommunicationCallMade
	CommunicationCallMergeIcon                      = icons.CommunicationCallMerge
	CommunicationCallMissedIcon                     = icons.CommunicationCallMissed
	CommunicationCallMissedOutgoingIcon             = icons.CommunicationCallMissedOutgoing
	CommunicationCallReceivedIcon                   = icons.CommunicationCallReceived
	CommunicationCallSplitIcon                      = icons.CommunicationCallSplit
	CommunicationChatIcon                           = icons.CommunicationChat
	CommunicationChatBubbleIcon                     = icons.CommunicationChatBubble
	CommunicationChatBubbleOutlineIcon              = icons.CommunicationChatBubbleOutline
	CommunicationClearAllIcon                       = icons.CommunicationClearAll
	CommunicationCommentIcon                        = icons.CommunicationComment
	CommunicationContactMailIcon                    = icons.CommunicationContactMail
	CommunicationContactPhoneIcon                   = icons.CommunicationContactPhone
	CommunicationContactsIcon                       = icons.CommunicationContacts
	CommunicationDialerSIPIcon                      = icons.CommunicationDialerSIP
	CommunicationDialpadIcon                        = icons.CommunicationDialpad
	CommunicationEmailIcon                          = icons.CommunicationEmail
	CommunicationForumIcon                          = icons.CommunicationForum
	CommunicationImportContactsIcon                 = icons.CommunicationImportContacts
	CommunicationImportExportIcon                   = icons.CommunicationImportExport
	CommunicationInvertColorsOffIcon                = icons.CommunicationInvertColorsOff
	CommunicationLiveHelpIcon                       = icons.CommunicationLiveHelp
	CommunicationLocationOffIcon                    = icons.CommunicationLocationOff
	CommunicationLocationOnIcon                     = icons.CommunicationLocationOn
	CommunicationMailOutlineIcon                    = icons.CommunicationMailOutline
	CommunicationMessageIcon                        = icons.CommunicationMessage
	CommunicationNoSIMIcon                          = icons.CommunicationNoSIM
	CommunicationPhoneIcon                          = icons.CommunicationPhone
	CommunicationPhoneLinkEraseIcon                 = icons.CommunicationPhoneLinkErase
	CommunicationPhoneLinkLockIcon                  = icons.CommunicationPhoneLinkLock
	CommunicationPhoneLinkRingIcon                  = icons.CommunicationPhoneLinkRing
	CommunicationPhoneLinkSetupIcon                 = icons.CommunicationPhoneLinkSetup
	CommunicationPortableWiFiOffIcon                = icons.CommunicationPortableWiFiOff
	CommunicationPresentToAllIcon                   = icons.CommunicationPresentToAll
	CommunicationRingVolumeIcon                     = icons.CommunicationRingVolume
	CommunicationRSSFeedIcon                        = icons.CommunicationRSSFeed
	CommunicationScreenShareIcon                    = icons.CommunicationScreenShare
	CommunicationSpeakerPhoneIcon                   = icons.CommunicationSpeakerPhone
	CommunicationStayCurrentLandscapeIcon           = icons.CommunicationStayCurrentLandscape
	CommunicationStayCurrentPortraitIcon            = icons.CommunicationStayCurrentPortrait
	CommunicationStayPrimaryLandscapeIcon           = icons.CommunicationStayPrimaryLandscape
	CommunicationStayPrimaryPortraitIcon            = icons.CommunicationStayPrimaryPortrait
	CommunicationStopScreenShareIcon                = icons.CommunicationStopScreenShare
	CommunicationSwapCallsIcon                      = icons.CommunicationSwapCalls
	CommunicationTextSMSIcon                        = icons.CommunicationTextSMS
	CommunicationVoicemailIcon                      = icons.CommunicationVoicemail
	CommunicationVPNKeyIcon                         = icons.CommunicationVPNKey
	ContentAddIcon                                  = icons.ContentAdd
	ContentAddBoxIcon                               = icons.ContentAddBox
	ContentAddCircleIcon                            = icons.ContentAddCircle
	ContentAddCircleOutlineIcon                     = icons.ContentAddCircleOutline
	ContentArchiveIcon                              = icons.ContentArchive
	ContentBackspaceIcon                            = icons.ContentBackspace
	ContentBlockIcon                                = icons.ContentBlock
	ContentClearIcon                                = icons.ContentClear
	ContentContentCopyIcon                          = icons.ContentContentCopy
	ContentContentCutIcon                           = icons.ContentContentCut
	ContentContentPasteIcon                         = icons.ContentContentPaste
	ContentCreateIcon                               = icons.ContentCreate
	ContentDeleteSweepIcon                          = icons.ContentDeleteSweep
	ContentDraftsIcon                               = icons.ContentDrafts
	ContentFilterListIcon                           = icons.ContentFilterList
	ContentFlagIcon                                 = icons.ContentFlag
	ContentFontDownloadIcon                         = icons.ContentFontDownload
	ContentForwardIcon                              = icons.ContentForward
	ContentGestureIcon                              = icons.ContentGesture
	ContentInboxIcon                                = icons.ContentInbox
	ContentLinkIcon                                 = icons.ContentLink
	ContentLowPriorityIcon                          = icons.ContentLowPriority
	ContentMailIcon                                 = icons.ContentMail
	ContentMarkUnreadIcon                           = icons.ContentMarkUnread
	ContentMoveToInboxIcon                          = icons.ContentMoveToInbox
	ContentNextWeekIcon                             = icons.ContentNextWeek
	ContentRedoIcon                                 = icons.ContentRedo
	ContentRemoveIcon                               = icons.ContentRemove
	ContentRemoveCircleIcon                         = icons.ContentRemoveCircle
	ContentRemoveCircleOutlineIcon                  = icons.ContentRemoveCircleOutline
	ContentReplyIcon                                = icons.ContentReply
	ContentReplyAllIcon                             = icons.ContentReplyAll
	ContentReportIcon                               = icons.ContentReport
	ContentSaveIcon                                 = icons.ContentSave
	ContentSelectAllIcon                            = icons.ContentSelectAll
	ContentSendIcon                                 = icons.ContentSend
	ContentSortIcon                                 = icons.ContentSort
	ContentTextFormatIcon                           = icons.ContentTextFormat
	ContentUnarchiveIcon                            = icons.ContentUnarchive
	ContentUndoIcon                                 = icons.ContentUndo
	ContentWeekendIcon                              = icons.ContentWeekend
	DeviceAccessAlarmIcon                           = icons.DeviceAccessAlarm
	DeviceAccessAlarmsIcon                          = icons.DeviceAccessAlarms
	DeviceAccessTimeIcon                            = icons.DeviceAccessTime
	DeviceAddAlarmIcon                              = icons.DeviceAddAlarm
	DeviceAirplaneModeActiveIcon                    = icons.DeviceAirplaneModeActive
	DeviceAirplaneModeInactiveIcon                  = icons.DeviceAirplaneModeInactive
	DeviceBattery20Icon                             = icons.DeviceBattery20
	DeviceBattery30Icon                             = icons.DeviceBattery30
	DeviceBattery50Icon                             = icons.DeviceBattery50
	DeviceBattery60Icon                             = icons.DeviceBattery60
	DeviceBattery80Icon                             = icons.DeviceBattery80
	DeviceBattery90Icon                             = icons.DeviceBattery90
	DeviceBatteryAlertIcon                          = icons.DeviceBatteryAlert
	DeviceBatteryCharging20Icon                     = icons.DeviceBatteryCharging20
	DeviceBatteryCharging30Icon                     = icons.DeviceBatteryCharging30
	DeviceBatteryCharging50Icon                     = icons.DeviceBatteryCharging50
	DeviceBatteryCharging60Icon                     = icons.DeviceBatteryCharging60
	DeviceBatteryCharging80Icon                     = icons.DeviceBatteryCharging80
	DeviceBatteryCharging90Icon                     = icons.DeviceBatteryCharging90
	DeviceBatteryChargingFullIcon                   = icons.DeviceBatteryChargingFull
	DeviceBatteryFullIcon                           = icons.DeviceBatteryFull
	DeviceBatteryStdIcon                            = icons.DeviceBatteryStd
	DeviceBatteryUnknownIcon                        = icons.DeviceBatteryUnknown
	DeviceBluetoothIcon                             = icons.DeviceBluetooth
	DeviceBluetoothConnectedIcon                    = icons.DeviceBluetoothConnected
	DeviceBluetoothDisabledIcon                     = icons.DeviceBluetoothDisabled
	DeviceBluetoothSearchingIcon                    = icons.DeviceBluetoothSearching
	DeviceBrightnessAutoIcon                        = icons.DeviceBrightnessAuto
	DeviceBrightnessHighIcon                        = icons.DeviceBrightnessHigh
	DeviceBrightnessLowIcon                         = icons.DeviceBrightnessLow
	DeviceBrightnessMediumIcon                      = icons.DeviceBrightnessMedium
	DeviceDataUsageIcon                             = icons.DeviceDataUsage
	DeviceDeveloperModeIcon                         = icons.DeviceDeveloperMode
	DeviceDevicesIcon                               = icons.DeviceDevices
	DeviceDVRIcon                                   = icons.DeviceDVR
	DeviceGPSFixedIcon                              = icons.DeviceGPSFixed
	DeviceGPSNotFixedIcon                           = icons.DeviceGPSNotFixed
	DeviceGPSOffIcon                                = icons.DeviceGPSOff
	DeviceGraphicEqIcon                             = icons.DeviceGraphicEq
	DeviceLocationDisabledIcon                      = icons.DeviceLocationDisabled
	DeviceLocationSearchingIcon                     = icons.DeviceLocationSearching
	DeviceNetworkCellIcon                           = icons.DeviceNetworkCell
	DeviceNetworkWiFiIcon                           = icons.DeviceNetworkWiFi
	DeviceNFCIcon                                   = icons.DeviceNFC
	DeviceScreenLockLandscapeIcon                   = icons.DeviceScreenLockLandscape
	DeviceScreenLockPortraitIcon                    = icons.DeviceScreenLockPortrait
	DeviceScreenLockRotationIcon                    = icons.DeviceScreenLockRotation
	DeviceScreenRotationIcon                        = icons.DeviceScreenRotation
	DeviceSDStorageIcon                             = icons.DeviceSDStorage
	DeviceSettingsSystemDaydreamIcon                = icons.DeviceSettingsSystemDaydream
	DeviceSignalCellular0BarIcon                    = icons.DeviceSignalCellular0Bar
	DeviceSignalCellular1BarIcon                    = icons.DeviceSignalCellular1Bar
	DeviceSignalCellular2BarIcon                    = icons.DeviceSignalCellular2Bar
	DeviceSignalCellular3BarIcon                    = icons.DeviceSignalCellular3Bar
	DeviceSignalCellular4BarIcon                    = icons.DeviceSignalCellular4Bar
	DeviceSignalCellularConnectedNoInternet0BarIcon = icons.DeviceSignalCellularConnectedNoInternet0Bar
	DeviceSignalCellularConnectedNoInternet1BarIcon = icons.DeviceSignalCellularConnectedNoInternet1Bar
	DeviceSignalCellularConnectedNoInternet2BarIcon = icons.DeviceSignalCellularConnectedNoInternet2Bar
	DeviceSignalCellularConnectedNoInternet3BarIcon = icons.DeviceSignalCellularConnectedNoInternet3Bar
	DeviceSignalCellularConnectedNoInternet4BarIcon = icons.DeviceSignalCellularConnectedNoInternet4Bar
	DeviceSignalCellularNoSIMIcon                   = icons.DeviceSignalCellularNoSIM
	DeviceSignalCellularNullIcon                    = icons.DeviceSignalCellularNull
	DeviceSignalCellularOffIcon                     = icons.DeviceSignalCellularOff
	DeviceSignalWiFi0BarIcon                        = icons.DeviceSignalWiFi0Bar
	DeviceSignalWiFi1BarIcon                        = icons.DeviceSignalWiFi1Bar
	DeviceSignalWiFi1BarLockIcon                    = icons.DeviceSignalWiFi1BarLock
	DeviceSignalWiFi2BarIcon                        = icons.DeviceSignalWiFi2Bar
	DeviceSignalWiFi2BarLockIcon                    = icons.DeviceSignalWiFi2BarLock
	DeviceSignalWiFi3BarIcon                        = icons.DeviceSignalWiFi3Bar
	DeviceSignalWiFi3BarLockIcon                    = icons.DeviceSignalWiFi3BarLock
	DeviceSignalWiFi4BarIcon                        = icons.DeviceSignalWiFi4Bar
	DeviceSignalWiFi4BarLockIcon                    = icons.DeviceSignalWiFi4BarLock
	DeviceSignalWiFiOffIcon                         = icons.DeviceSignalWiFiOff
	DeviceStorageIcon                               = icons.DeviceStorage
	DeviceUSBIcon                                   = icons.DeviceUSB
	DeviceWallpaperIcon                             = icons.DeviceWallpaper
	DeviceWidgetsIcon                               = icons.DeviceWidgets
	DeviceWiFiLockIcon                              = icons.DeviceWiFiLock
	DeviceWiFiTetheringIcon                         = icons.DeviceWiFiTethering
	EditorAttachFileIcon                            = icons.EditorAttachFile
	EditorAttachMoneyIcon                           = icons.EditorAttachMoney
	EditorBorderAllIcon                             = icons.EditorBorderAll
	EditorBorderBottomIcon                          = icons.EditorBorderBottom
	EditorBorderClearIcon                           = icons.EditorBorderClear
	EditorBorderColorIcon                           = icons.EditorBorderColor
	EditorBorderHorizontalIcon                      = icons.EditorBorderHorizontal
	EditorBorderInnerIcon                           = icons.EditorBorderInner
	EditorBorderLeftIcon                            = icons.EditorBorderLeft
	EditorBorderOuterIcon                           = icons.EditorBorderOuter
	EditorBorderRightIcon                           = icons.EditorBorderRight
	EditorBorderStyleIcon                           = icons.EditorBorderStyle
	EditorBorderTopIcon                             = icons.EditorBorderTop
	EditorBorderVerticalIcon                        = icons.EditorBorderVertical
	EditorBubbleChartIcon                           = icons.EditorBubbleChart
	EditorDragHandleIcon                            = icons.EditorDragHandle
	EditorFormatAlignCenterIcon                     = icons.EditorFormatAlignCenter
	EditorFormatAlignJustifyIcon                    = icons.EditorFormatAlignJustify
	EditorFormatAlignLeftIcon                       = icons.EditorFormatAlignLeft
	EditorFormatAlignRightIcon                      = icons.EditorFormatAlignRight
	EditorFormatBoldIcon                            = icons.EditorFormatBold
	EditorFormatClearIcon                           = icons.EditorFormatClear
	EditorFormatColorFillIcon                       = icons.EditorFormatColorFill
	EditorFormatColorResetIcon                      = icons.EditorFormatColorReset
	EditorFormatColorTextIcon                       = icons.EditorFormatColorText
	EditorFormatIndentDecreaseIcon                  = icons.EditorFormatIndentDecrease
	EditorFormatIndentIncreaseIcon                  = icons.EditorFormatIndentIncrease
	EditorFormatItalicIcon                          = icons.EditorFormatItalic
	EditorFormatLineSpacingIcon                     = icons.EditorFormatLineSpacing
	EditorFormatListBulletedIcon                    = icons.EditorFormatListBulleted
	EditorFormatListNumberedIcon                    = icons.EditorFormatListNumbered
	EditorFormatPaintIcon                           = icons.EditorFormatPaint
	EditorFormatQuoteIcon                           = icons.EditorFormatQuote
	EditorFormatShapesIcon                          = icons.EditorFormatShapes
	EditorFormatSizeIcon                            = icons.EditorFormatSize
	EditorFormatStrikethroughIcon                   = icons.EditorFormatStrikethrough
	EditorFormatTextDirectionLToRIcon               = icons.EditorFormatTextDirectionLToR
	EditorFormatTextDirectionRToLIcon               = icons.EditorFormatTextDirectionRToL
	EditorFormatUnderlinedIcon                      = icons.EditorFormatUnderlined
	EditorFunctionsIcon                             = icons.EditorFunctions
	EditorHighlightIcon                             = icons.EditorHighlight
	EditorInsertChartIcon                           = icons.EditorInsertChart
	EditorInsertCommentIcon                         = icons.EditorInsertComment
	EditorInsertDriveFileIcon                       = icons.EditorInsertDriveFile
	EditorInsertEmoticonIcon                        = icons.EditorInsertEmoticon
	EditorInsertInvitationIcon                      = icons.EditorInsertInvitation
	EditorInsertLinkIcon                            = icons.EditorInsertLink
	EditorInsertPhotoIcon                           = icons.EditorInsertPhoto
	EditorLinearScaleIcon                           = icons.EditorLinearScale
	EditorMergeTypeIcon                             = icons.EditorMergeType
	EditorModeCommentIcon                           = icons.EditorModeComment
	EditorModeEditIcon                              = icons.EditorModeEdit
	EditorMonetizationOnIcon                        = icons.EditorMonetizationOn
	EditorMoneyOffIcon                              = icons.EditorMoneyOff
	EditorMultilineChartIcon                        = icons.EditorMultilineChart
	EditorPieChartIcon                              = icons.EditorPieChart
	EditorPieChartOutlinedIcon                      = icons.EditorPieChartOutlined
	EditorPublishIcon                               = icons.EditorPublish
	EditorShortTextIcon                             = icons.EditorShortText
	EditorShowChartIcon                             = icons.EditorShowChart
	EditorSpaceBarIcon                              = icons.EditorSpaceBar
	EditorStrikethroughSIcon                        = icons.EditorStrikethroughS
	EditorTextFieldsIcon                            = icons.EditorTextFields
	EditorTitleIcon                                 = icons.EditorTitle
	EditorVerticalAlignBottomIcon                   = icons.EditorVerticalAlignBottom
	EditorVerticalAlignCenterIcon                   = icons.EditorVerticalAlignCenter
	EditorVerticalAlignTopIcon                      = icons.EditorVerticalAlignTop
	EditorWrapTextIcon                              = icons.EditorWrapText
	FileAttachmentIcon                              = icons.FileAttachment
	FileCloudIcon                                   = icons.FileCloud
	FileCloudCircleIcon                             = icons.FileCloudCircle
	FileCloudDoneIcon                               = icons.FileCloudDone
	FileCloudDownloadIcon                           = icons.FileCloudDownload
	FileCloudOffIcon                                = icons.FileCloudOff
	FileCloudQueueIcon                              = icons.FileCloudQueue
	FileCloudUploadIcon                             = icons.FileCloudUpload
	FileCreateNewFolderIcon                         = icons.FileCreateNewFolder
	FileFileDownloadIcon                            = icons.FileFileDownload
	FileFileUploadIcon                              = icons.FileFileUpload
	FileFolderIcon                                  = icons.FileFolder
	FileFolderOpenIcon                              = icons.FileFolderOpen
	FileFolderSharedIcon                            = icons.FileFolderShared
	HardwareCastIcon                                = icons.HardwareCast
	HardwareCastConnectedIcon                       = icons.HardwareCastConnected
	HardwareComputerIcon                            = icons.HardwareComputer
	HardwareDesktopMacIcon                          = icons.HardwareDesktopMac
	HardwareDesktopWindowsIcon                      = icons.HardwareDesktopWindows
	HardwareDeveloperBoardIcon                      = icons.HardwareDeveloperBoard
	HardwareDeviceHubIcon                           = icons.HardwareDeviceHub
	HardwareDevicesOtherIcon                        = icons.HardwareDevicesOther
	HardwareDockIcon                                = icons.HardwareDock
	HardwareGamepadIcon                             = icons.HardwareGamepad
	HardwareHeadsetIcon                             = icons.HardwareHeadset
	HardwareHeadsetMicIcon                          = icons.HardwareHeadsetMic
	HardwareKeyboardIcon                            = icons.HardwareKeyboard
	HardwareKeyboardArrowDownIcon                   = icons.HardwareKeyboardArrowDown
	HardwareKeyboardArrowLeftIcon                   = icons.HardwareKeyboardArrowLeft
	HardwareKeyboardArrowRightIcon                  = icons.HardwareKeyboardArrowRight
	HardwareKeyboardArrowUpIcon                     = icons.HardwareKeyboardArrowUp
	HardwareKeyboardBackspaceIcon                   = icons.HardwareKeyboardBackspace
	HardwareKeyboardCapslockIcon                    = icons.HardwareKeyboardCapslock
	HardwareKeyboardHideIcon                        = icons.HardwareKeyboardHide
	HardwareKeyboardReturnIcon                      = icons.HardwareKeyboardReturn
	HardwareKeyboardTabIcon                         = icons.HardwareKeyboardTab
	HardwareKeyboardVoiceIcon                       = icons.HardwareKeyboardVoice
	HardwareLaptopIcon                              = icons.HardwareLaptop
	HardwareLaptopChromebookIcon                    = icons.HardwareLaptopChromebook
	HardwareLaptopMacIcon                           = icons.HardwareLaptopMac
	HardwareLaptopWindowsIcon                       = icons.HardwareLaptopWindows
	HardwareMemoryIcon                              = icons.HardwareMemory
	HardwareMouseIcon                               = icons.HardwareMouse
	HardwarePhoneAndroidIcon                        = icons.HardwarePhoneAndroid
	HardwarePhoneIPhoneIcon                         = icons.HardwarePhoneIPhone
	HardwarePhoneLinkIcon                           = icons.HardwarePhoneLink
	HardwarePhoneLinkOffIcon                        = icons.HardwarePhoneLinkOff
	HardwarePowerInputIcon                          = icons.HardwarePowerInput
	HardwareRouterIcon                              = icons.HardwareRouter
	HardwareScannerIcon                             = icons.HardwareScanner
	HardwareSecurityIcon                            = icons.HardwareSecurity
	HardwareSIMCardIcon                             = icons.HardwareSIMCard
	HardwareSmartphoneIcon                          = icons.HardwareSmartphone
	HardwareSpeakerIcon                             = icons.HardwareSpeaker
	HardwareSpeakerGroupIcon                        = icons.HardwareSpeakerGroup
	HardwareTabletIcon                              = icons.HardwareTablet
	HardwareTabletAndroidIcon                       = icons.HardwareTabletAndroid
	HardwareTabletMacIcon                           = icons.HardwareTabletMac
	HardwareToysIcon                                = icons.HardwareToys
	HardwareTVIcon                                  = icons.HardwareTV
	HardwareVideogameAssetIcon                      = icons.HardwareVideogameAsset
	HardwareWatchIcon                               = icons.HardwareWatch
	ImageAddAPhotoIcon                              = icons.ImageAddAPhoto
	ImageAddToPhotosIcon                            = icons.ImageAddToPhotos
	ImageAdjustIcon                                 = icons.ImageAdjust
	ImageAssistantIcon                              = icons.ImageAssistant
	ImageAssistantPhotoIcon                         = icons.ImageAssistantPhoto
	ImageAudiotrackIcon                             = icons.ImageAudiotrack
	ImageBlurCircularIcon                           = icons.ImageBlurCircular
	ImageBlurLinearIcon                             = icons.ImageBlurLinear
	ImageBlurOffIcon                                = icons.ImageBlurOff
	ImageBlurOnIcon                                 = icons.ImageBlurOn
	ImageBrightness1Icon                            = icons.ImageBrightness1
	ImageBrightness2Icon                            = icons.ImageBrightness2
	ImageBrightness3Icon                            = icons.ImageBrightness3
	ImageBrightness4Icon                            = icons.ImageBrightness4
	ImageBrightness5Icon                            = icons.ImageBrightness5
	ImageBrightness6Icon                            = icons.ImageBrightness6
	ImageBrightness7Icon                            = icons.ImageBrightness7
	ImageBrokenImageIcon                            = icons.ImageBrokenImage
	ImageBrushIcon                                  = icons.ImageBrush
	ImageBurstModeIcon                              = icons.ImageBurstMode
	ImageCameraIcon                                 = icons.ImageCamera
	ImageCameraAltIcon                              = icons.ImageCameraAlt
	ImageCameraFrontIcon                            = icons.ImageCameraFront
	ImageCameraRearIcon                             = icons.ImageCameraRear
	ImageCameraRollIcon                             = icons.ImageCameraRoll
	ImageCenterFocusStrongIcon                      = icons.ImageCenterFocusStrong
	ImageCenterFocusWeakIcon                        = icons.ImageCenterFocusWeak
	ImageCollectionsIcon                            = icons.ImageCollections
	ImageCollectionsBookmarkIcon                    = icons.ImageCollectionsBookmark
	ImageColorLensIcon                              = icons.ImageColorLens
	ImageColorizeIcon                               = icons.ImageColorize
	ImageCompareIcon                                = icons.ImageCompare
	ImageControlPointIcon                           = icons.ImageControlPoint
	ImageControlPointDuplicateIcon                  = icons.ImageControlPointDuplicate
	ImageCropIcon                                   = icons.ImageCrop
	ImageCrop169Icon                                = icons.ImageCrop169
	ImageCrop32Icon                                 = icons.ImageCrop32
	ImageCrop54Icon                                 = icons.ImageCrop54
	ImageCrop75Icon                                 = icons.ImageCrop75
	ImageCropDINIcon                                = icons.ImageCropDIN
	ImageCropFreeIcon                               = icons.ImageCropFree
	ImageCropLandscapeIcon                          = icons.ImageCropLandscape
	ImageCropOriginalIcon                           = icons.ImageCropOriginal
	ImageCropPortraitIcon                           = icons.ImageCropPortrait
	ImageCropRotateIcon                             = icons.ImageCropRotate
	ImageCropSquareIcon                             = icons.ImageCropSquare
	ImageDehazeIcon                                 = icons.ImageDehaze
	ImageDetailsIcon                                = icons.ImageDetails
	ImageEditIcon                                   = icons.ImageEdit
	ImageExposureIcon                               = icons.ImageExposure
	ImageExposureNeg1Icon                           = icons.ImageExposureNeg1
	ImageExposureNeg2Icon                           = icons.ImageExposureNeg2
	ImageExposurePlus1Icon                          = icons.ImageExposurePlus1
	ImageExposurePlus2Icon                          = icons.ImageExposurePlus2
	ImageExposureZeroIcon                           = icons.ImageExposureZero
	ImageFilterIcon                                 = icons.ImageFilter
	ImageFilter1Icon                                = icons.ImageFilter1
	ImageFilter2Icon                                = icons.ImageFilter2
	ImageFilter3Icon                                = icons.ImageFilter3
	ImageFilter4Icon                                = icons.ImageFilter4
	ImageFilter5Icon                                = icons.ImageFilter5
	ImageFilter6Icon                                = icons.ImageFilter6
	ImageFilter7Icon                                = icons.ImageFilter7
	ImageFilter8Icon                                = icons.ImageFilter8
	ImageFilter9Icon                                = icons.ImageFilter9
	ImageFilter9PlusIcon                            = icons.ImageFilter9Plus
	ImageFilterBAndWIcon                            = icons.ImageFilterBAndW
	ImageFilterCenterFocusIcon                      = icons.ImageFilterCenterFocus
	ImageFilterDramaIcon                            = icons.ImageFilterDrama
	ImageFilterFramesIcon                           = icons.ImageFilterFrames
	ImageFilterHDRIcon                              = icons.ImageFilterHDR
	ImageFilterNoneIcon                             = icons.ImageFilterNone
	ImageFilterTiltShiftIcon                        = icons.ImageFilterTiltShift
	ImageFilterVintageIcon                          = icons.ImageFilterVintage
	ImageFlareIcon                                  = icons.ImageFlare
	ImageFlashAutoIcon                              = icons.ImageFlashAuto
	ImageFlashOffIcon                               = icons.ImageFlashOff
	ImageFlashOnIcon                                = icons.ImageFlashOn
	ImageFlipIcon                                   = icons.ImageFlip
	ImageGradientIcon                               = icons.ImageGradient
	ImageGrainIcon                                  = icons.ImageGrain
	ImageGridOffIcon                                = icons.ImageGridOff
	ImageGridOnIcon                                 = icons.ImageGridOn
	ImageHDROffIcon                                 = icons.ImageHDROff
	ImageHDROnIcon                                  = icons.ImageHDROn
	ImageHDRStrongIcon                              = icons.ImageHDRStrong
	ImageHDRWeakIcon                                = icons.ImageHDRWeak
	ImageHealingIcon                                = icons.ImageHealing
	ImageImageIcon                                  = icons.ImageImage
	ImageImageAspectRatioIcon                       = icons.ImageImageAspectRatio
	ImageISOIcon                                    = icons.ImageISO
	ImageLandscapeIcon                              = icons.ImageLandscape
	ImageLeakAddIcon                                = icons.ImageLeakAdd
	ImageLeakRemoveIcon                             = icons.ImageLeakRemove
	ImageLensIcon                                   = icons.ImageLens
	ImageLinkedCameraIcon                           = icons.ImageLinkedCamera
	ImageLooksIcon                                  = icons.ImageLooks
	ImageLooks3Icon                                 = icons.ImageLooks3
	ImageLooks4Icon                                 = icons.ImageLooks4
	ImageLooks5Icon                                 = icons.ImageLooks5
	ImageLooks6Icon                                 = icons.ImageLooks6
	ImageLooksOneIcon                               = icons.ImageLooksOne
	ImageLooksTwoIcon                               = icons.ImageLooksTwo
	ImageLoupeIcon                                  = icons.ImageLoupe
	ImageMonochromePhotosIcon                       = icons.ImageMonochromePhotos
	ImageMovieCreationIcon                          = icons.ImageMovieCreation
	ImageMovieFilterIcon                            = icons.ImageMovieFilter
	ImageMusicNoteIcon                              = icons.ImageMusicNote
	ImageNatureIcon                                 = icons.ImageNature
	ImageNaturePeopleIcon                           = icons.ImageNaturePeople
	ImageNavigateBeforeIcon                         = icons.ImageNavigateBefore
	ImageNavigateNextIcon                           = icons.ImageNavigateNext
	ImagePaletteIcon                                = icons.ImagePalette
	ImagePanoramaIcon                               = icons.ImagePanorama
	ImagePanoramaFishEyeIcon                        = icons.ImagePanoramaFishEye
	ImagePanoramaHorizontalIcon                     = icons.ImagePanoramaHorizontal
	ImagePanoramaVerticalIcon                       = icons.ImagePanoramaVertical
	ImagePanoramaWideAngleIcon                      = icons.ImagePanoramaWideAngle
	ImagePhotoIcon                                  = icons.ImagePhoto
	ImagePhotoAlbumIcon                             = icons.ImagePhotoAlbum
	ImagePhotoCameraIcon                            = icons.ImagePhotoCamera
	ImagePhotoFilterIcon                            = icons.ImagePhotoFilter
	ImagePhotoLibraryIcon                           = icons.ImagePhotoLibrary
	ImagePhotoSizeSelectActualIcon                  = icons.ImagePhotoSizeSelectActual
	ImagePhotoSizeSelectLargeIcon                   = icons.ImagePhotoSizeSelectLarge
	ImagePhotoSizeSelectSmallIcon                   = icons.ImagePhotoSizeSelectSmall
	ImagePictureAsPDFIcon                           = icons.ImagePictureAsPDF
	ImagePortraitIcon                               = icons.ImagePortrait
	ImageRemoveRedEyeIcon                           = icons.ImageRemoveRedEye
	ImageRotate90DegreesCCWIcon                     = icons.ImageRotate90DegreesCCW
	ImageRotateLeftIcon                             = icons.ImageRotateLeft
	ImageRotateRightIcon                            = icons.ImageRotateRight
	ImageSlideshowIcon                              = icons.ImageSlideshow
	ImageStraightenIcon                             = icons.ImageStraighten
	ImageStyleIcon                                  = icons.ImageStyle
	ImageSwitchCameraIcon                           = icons.ImageSwitchCamera
	ImageSwitchVideoIcon                            = icons.ImageSwitchVideo
	ImageTagFacesIcon                               = icons.ImageTagFaces
	ImageTextureIcon                                = icons.ImageTexture
	ImageTimeLapseIcon                              = icons.ImageTimeLapse
	ImageTimerIcon                                  = icons.ImageTimer
	ImageTimer10Icon                                = icons.ImageTimer10
	ImageTimer3Icon                                 = icons.ImageTimer3
	ImageTimerOffIcon                               = icons.ImageTimerOff
	ImageTonalityIcon                               = icons.ImageTonality
	ImageTransformIcon                              = icons.ImageTransform
	ImageTuneIcon                                   = icons.ImageTune
	ImageViewComfyIcon                              = icons.ImageViewComfy
	ImageViewCompactIcon                            = icons.ImageViewCompact
	ImageVignetteIcon                               = icons.ImageVignette
	ImageWBAutoIcon                                 = icons.ImageWBAuto
	ImageWBCloudyIcon                               = icons.ImageWBCloudy
	ImageWBIncandescentIcon                         = icons.ImageWBIncandescent
	ImageWBIridescentIcon                           = icons.ImageWBIridescent
	ImageWBSunnyIcon                                = icons.ImageWBSunny
	MapsAddLocationIcon                             = icons.MapsAddLocation
	MapsBeenhereIcon                                = icons.MapsBeenhere
	MapsDirectionsIcon                              = icons.MapsDirections
	MapsDirectionsBikeIcon                          = icons.MapsDirectionsBike
	MapsDirectionsBoatIcon                          = icons.MapsDirectionsBoat
	MapsDirectionsBusIcon                           = icons.MapsDirectionsBus
	MapsDirectionsCarIcon                           = icons.MapsDirectionsCar
	MapsDirectionsRailwayIcon                       = icons.MapsDirectionsRailway
	MapsDirectionsRunIcon                           = icons.MapsDirectionsRun
	MapsDirectionsSubwayIcon                        = icons.MapsDirectionsSubway
	MapsDirectionsTransitIcon                       = icons.MapsDirectionsTransit
	MapsDirectionsWalkIcon                          = icons.MapsDirectionsWalk
	MapsEditLocationIcon                            = icons.MapsEditLocation
	MapsEVStationIcon                               = icons.MapsEVStation
	MapsFlightIcon                                  = icons.MapsFlight
	MapsHotelIcon                                   = icons.MapsHotel
	MapsLayersIcon                                  = icons.MapsLayers
	MapsLayersClearIcon                             = icons.MapsLayersClear
	MapsLocalActivityIcon                           = icons.MapsLocalActivity
	MapsLocalAirportIcon                            = icons.MapsLocalAirport
	MapsLocalATMIcon                                = icons.MapsLocalATM
	MapsLocalBarIcon                                = icons.MapsLocalBar
	MapsLocalCafeIcon                               = icons.MapsLocalCafe
	MapsLocalCarWashIcon                            = icons.MapsLocalCarWash
	MapsLocalConvenienceStoreIcon                   = icons.MapsLocalConvenienceStore
	MapsLocalDiningIcon                             = icons.MapsLocalDining
	MapsLocalDrinkIcon                              = icons.MapsLocalDrink
	MapsLocalFloristIcon                            = icons.MapsLocalFlorist
	MapsLocalGasStationIcon                         = icons.MapsLocalGasStation
	MapsLocalGroceryStoreIcon                       = icons.MapsLocalGroceryStore
	MapsLocalHospitalIcon                           = icons.MapsLocalHospital
	MapsLocalHotelIcon                              = icons.MapsLocalHotel
	MapsLocalLaundryServiceIcon                     = icons.MapsLocalLaundryService
	MapsLocalLibraryIcon                            = icons.MapsLocalLibrary
	MapsLocalMallIcon                               = icons.MapsLocalMall
	MapsLocalMoviesIcon                             = icons.MapsLocalMovies
	MapsLocalOfferIcon                              = icons.MapsLocalOffer
	MapsLocalParkingIcon                            = icons.MapsLocalParking
	MapsLocalPharmacyIcon                           = icons.MapsLocalPharmacy
	MapsLocalPhoneIcon                              = icons.MapsLocalPhone
	MapsLocalPizzaIcon                              = icons.MapsLocalPizza
	MapsLocalPlayIcon                               = icons.MapsLocalPlay
	MapsLocalPostOfficeIcon                         = icons.MapsLocalPostOffice
	MapsLocalPrintshopIcon                          = icons.MapsLocalPrintshop
	MapsLocalSeeIcon                                = icons.MapsLocalSee
	MapsLocalShippingIcon                           = icons.MapsLocalShipping
	MapsLocalTaxiIcon                               = icons.MapsLocalTaxi
	MapsMapIcon                                     = icons.MapsMap
	MapsMyLocationIcon                              = icons.MapsMyLocation
	MapsNavigationIcon                              = icons.MapsNavigation
	MapsNearMeIcon                                  = icons.MapsNearMe
	MapsPersonPinIcon                               = icons.MapsPersonPin
	MapsPersonPinCircleIcon                         = icons.MapsPersonPinCircle
	MapsPinDropIcon                                 = icons.MapsPinDrop
	MapsPlaceIcon                                   = icons.MapsPlace
	MapsRateReviewIcon                              = icons.MapsRateReview
	MapsRestaurantIcon                              = icons.MapsRestaurant
	MapsRestaurantMenuIcon                          = icons.MapsRestaurantMenu
	MapsSatelliteIcon                               = icons.MapsSatellite
	MapsStoreMallDirectoryIcon                      = icons.MapsStoreMallDirectory
	MapsStreetViewIcon                              = icons.MapsStreetView
	MapsSubwayIcon                                  = icons.MapsSubway
	MapsTerrainIcon                                 = icons.MapsTerrain
	MapsTrafficIcon                                 = icons.MapsTraffic
	MapsTrainIcon                                   = icons.MapsTrain
	MapsTramIcon                                    = icons.MapsTram
	MapsTransferWithinAStationIcon                  = icons.MapsTransferWithinAStation
	MapsZoomOutMapIcon                              = icons.MapsZoomOutMap
	NavigationAppsIcon                              = icons.NavigationApps
	NavigationArrowBackIcon                         = icons.NavigationArrowBack
	NavigationArrowDownwardIcon                     = icons.NavigationArrowDownward
	NavigationArrowDropDownIcon                     = icons.NavigationArrowDropDown
	NavigationArrowDropDownCircleIcon               = icons.NavigationArrowDropDownCircle
	NavigationArrowDropUpIcon                       = icons.NavigationArrowDropUp
	NavigationArrowForwardIcon                      = icons.NavigationArrowForward
	NavigationArrowUpwardIcon                       = icons.NavigationArrowUpward
	NavigationCancelIcon                            = icons.NavigationCancel
	NavigationCheckIcon                             = icons.NavigationCheck
	NavigationChevronLeftIcon                       = icons.NavigationChevronLeft
	NavigationChevronRightIcon                      = icons.NavigationChevronRight
	NavigationCloseIcon                             = icons.NavigationClose
	NavigationExpandLessIcon                        = icons.NavigationExpandLess
	NavigationExpandMoreIcon                        = icons.NavigationExpandMore
	NavigationFirstPageIcon                         = icons.NavigationFirstPage
	NavigationFullscreenIcon                        = icons.NavigationFullscreen
	NavigationFullscreenExitIcon                    = icons.NavigationFullscreenExit
	NavigationLastPageIcon                          = icons.NavigationLastPage
	NavigationMenuIcon                              = icons.NavigationMenu
	NavigationMoreHorizIcon                         = icons.NavigationMoreHoriz
	NavigationMoreVertIcon                          = icons.NavigationMoreVert
	NavigationRefreshIcon                           = icons.NavigationRefresh
	NavigationSubdirectoryArrowLeftIcon             = icons.NavigationSubdirectoryArrowLeft
	NavigationSubdirectoryArrowRightIcon            = icons.NavigationSubdirectoryArrowRight
	NavigationUnfoldLessIcon                        = icons.NavigationUnfoldLess
	NavigationUnfoldMoreIcon                        = icons.NavigationUnfoldMore
	NotificationADBIcon                             = icons.NotificationADB
	NotificationAirlineSeatFlatIcon                 = icons.NotificationAirlineSeatFlat
	NotificationAirlineSeatFlatAngledIcon           = icons.NotificationAirlineSeatFlatAngled
	NotificationAirlineSeatIndividualSuiteIcon      = icons.NotificationAirlineSeatIndividualSuite
	NotificationAirlineSeatLegroomExtraIcon         = icons.NotificationAirlineSeatLegroomExtra
	NotificationAirlineSeatLegroomNormalIcon        = icons.NotificationAirlineSeatLegroomNormal
	NotificationAirlineSeatLegroomReducedIcon       = icons.NotificationAirlineSeatLegroomReduced
	NotificationAirlineSeatReclineExtraIcon         = icons.NotificationAirlineSeatReclineExtra
	NotificationAirlineSeatReclineNormalIcon        = icons.NotificationAirlineSeatReclineNormal
	NotificationBluetoothAudioIcon                  = icons.NotificationBluetoothAudio
	NotificationConfirmationNumberIcon              = icons.NotificationConfirmationNumber
	NotificationDiscFullIcon                        = icons.NotificationDiscFull
	NotificationDoNotDisturbIcon                    = icons.NotificationDoNotDisturb
	NotificationDoNotDisturbAltIcon                 = icons.NotificationDoNotDisturbAlt
	NotificationDoNotDisturbOffIcon                 = icons.NotificationDoNotDisturbOff
	NotificationDoNotDisturbOnIcon                  = icons.NotificationDoNotDisturbOn
	NotificationDriveETAIcon                        = icons.NotificationDriveETA
	NotificationEnhancedEncryptionIcon              = icons.NotificationEnhancedEncryption
	NotificationEventAvailableIcon                  = icons.NotificationEventAvailable
	NotificationEventBusyIcon                       = icons.NotificationEventBusy
	NotificationEventNoteIcon                       = icons.NotificationEventNote
	NotificationFolderSpecialIcon                   = icons.NotificationFolderSpecial
	NotificationLiveTVIcon                          = icons.NotificationLiveTV
	NotificationMMSIcon                             = icons.NotificationMMS
	NotificationMoreIcon                            = icons.NotificationMore
	NotificationNetworkCheckIcon                    = icons.NotificationNetworkCheck
	NotificationNetworkLockedIcon                   = icons.NotificationNetworkLocked
	NotificationNoEncryptionIcon                    = icons.NotificationNoEncryption
	NotificationOnDemandVideoIcon                   = icons.NotificationOnDemandVideo
	NotificationPersonalVideoIcon                   = icons.NotificationPersonalVideo
	NotificationPhoneBluetoothSpeakerIcon           = icons.NotificationPhoneBluetoothSpeaker
	NotificationPhoneForwardedIcon                  = icons.NotificationPhoneForwarded
	NotificationPhoneInTalkIcon                     = icons.NotificationPhoneInTalk
	NotificationPhoneLockedIcon                     = icons.NotificationPhoneLocked
	NotificationPhoneMissedIcon                     = icons.NotificationPhoneMissed
	NotificationPhonePausedIcon                     = icons.NotificationPhonePaused
	NotificationPowerIcon                           = icons.NotificationPower
	NotificationPriorityHighIcon                    = icons.NotificationPriorityHigh
	NotificationRVHookupIcon                        = icons.NotificationRVHookup
	NotificationSDCardIcon                          = icons.NotificationSDCard
	NotificationSIMCardAlertIcon                    = icons.NotificationSIMCardAlert
	NotificationSMSIcon                             = icons.NotificationSMS
	NotificationSMSFailedIcon                       = icons.NotificationSMSFailed
	NotificationSyncIcon                            = icons.NotificationSync
	NotificationSyncDisabledIcon                    = icons.NotificationSyncDisabled
	NotificationSyncProblemIcon                     = icons.NotificationSyncProblem
	NotificationSystemUpdateIcon                    = icons.NotificationSystemUpdate
	NotificationTapAndPlayIcon                      = icons.NotificationTapAndPlay
	NotificationTimeToLeaveIcon                     = icons.NotificationTimeToLeave
	NotificationVibrationIcon                       = icons.NotificationVibration
	NotificationVoiceChatIcon                       = icons.NotificationVoiceChat
	NotificationVPNLockIcon                         = icons.NotificationVPNLock
	NotificationWCIcon                              = icons.NotificationWC
	NotificationWiFiIcon                            = icons.NotificationWiFi
	PlacesACUnitIcon                                = icons.PlacesACUnit
	PlacesAirportShuttleIcon                        = icons.PlacesAirportShuttle
	PlacesAllInclusiveIcon                          = icons.PlacesAllInclusive
	PlacesBeachAccessIcon                           = icons.PlacesBeachAccess
	PlacesBusinessCenterIcon                        = icons.PlacesBusinessCenter
	PlacesCasinoIcon                                = icons.PlacesCasino
	PlacesChildCareIcon                             = icons.PlacesChildCare
	PlacesChildFriendlyIcon                         = icons.PlacesChildFriendly
	PlacesFitnessCenterIcon                         = icons.PlacesFitnessCenter
	PlacesFreeBreakfastIcon                         = icons.PlacesFreeBreakfast
	PlacesGolfCourseIcon                            = icons.PlacesGolfCourse
	PlacesHotTubIcon                                = icons.PlacesHotTub
	PlacesKitchenIcon                               = icons.PlacesKitchen
	PlacesPoolIcon                                  = icons.PlacesPool
	PlacesRoomServiceIcon                           = icons.PlacesRoomService
	PlacesRVHookupIcon                              = icons.PlacesRVHookup
	PlacesSmokeFreeIcon                             = icons.PlacesSmokeFree
	PlacesSmokingRoomsIcon                          = icons.PlacesSmokingRooms
	PlacesSpaIcon                                   = icons.PlacesSpa
	SocialCakeIcon                                  = icons.SocialCake
	SocialDomainIcon                                = icons.SocialDomain
	SocialGroupIcon                                 = icons.SocialGroup
	SocialGroupAddIcon                              = icons.SocialGroupAdd
	SocialLocationCityIcon                          = icons.SocialLocationCity
	SocialMoodIcon                                  = icons.SocialMood
	SocialMoodBadIcon                               = icons.SocialMoodBad
	SocialNotificationsIcon                         = icons.SocialNotifications
	SocialNotificationsActiveIcon                   = icons.SocialNotificationsActive
	SocialNotificationsNoneIcon                     = icons.SocialNotificationsNone
	SocialNotificationsOffIcon                      = icons.SocialNotificationsOff
	SocialNotificationsPausedIcon                   = icons.SocialNotificationsPaused
	SocialPagesIcon                                 = icons.SocialPages
	SocialPartyModeIcon                             = icons.SocialPartyMode
	SocialPeopleIcon                                = icons.SocialPeople
	SocialPeopleOutlineIcon                         = icons.SocialPeopleOutline
	SocialPersonIcon                                = icons.SocialPerson
	SocialPersonAddIcon                             = icons.SocialPersonAdd
	SocialPersonOutlineIcon                         = icons.SocialPersonOutline
	SocialPlusOneIcon                               = icons.SocialPlusOne
	SocialPollIcon                                  = icons.SocialPoll
	SocialPublicIcon                                = icons.SocialPublic
	SocialSchoolIcon                                = icons.SocialSchool
	SocialSentimentDissatisfiedIcon                 = icons.SocialSentimentDissatisfied
	SocialSentimentNeutralIcon                      = icons.SocialSentimentNeutral
	SocialSentimentSatisfiedIcon                    = icons.SocialSentimentSatisfied
	SocialSentimentVeryDissatisfiedIcon             = icons.SocialSentimentVeryDissatisfied
	SocialSentimentVerySatisfiedIcon                = icons.SocialSentimentVerySatisfied
	SocialShareIcon                                 = icons.SocialShare
	SocialWhatsHotIcon                              = icons.SocialWhatsHot
	ToggleCheckBoxIcon                              = icons.ToggleCheckBox
	ToggleCheckBoxOutlineBlankIcon                  = icons.ToggleCheckBoxOutlineBlank
	ToggleIndeterminateCheckBoxIcon                 = icons.ToggleIndeterminateCheckBox
	ToggleRadioButtonCheckedIcon                    = icons.ToggleRadioButtonChecked
	ToggleRadioButtonUncheckedIcon                  = icons.ToggleRadioButtonUnchecked
	ToggleStarIcon                                  = icons.ToggleStar
	ToggleStarBorderIcon                            = icons.ToggleStarBorder
	ToggleStarHalfIcon                              = icons.ToggleStarHalf
)
