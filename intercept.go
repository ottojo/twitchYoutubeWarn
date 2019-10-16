package main

import (
	"github.com/gempir/go-twitch-irc"
	"log"
	"regexp"
	"strings"
)

// intercept looks for a youtube link / video id and issues a warning if appropriate
func intercept(message twitch.PrivateMessage, matchedPrefix string) {

	reason, suspicious := analyzeMessage(message)

	if suspicious {
		log.Printf("Intercept | Suspicious | %s | %s | \"%s\" | %s\n",
			message.Channel, message.User.DisplayName, message.Message, reason)
		twitchChatClient.Say(message.Channel, "Warning: "+reason)
	} else {
		log.Printf("Intercept | Not Suspicious | %s | %s | \"%s\"\n",
			message.Channel, message.User.DisplayName, message.Message)
	}
}

func analyzeMessage(message twitch.PrivateMessage) (reason string, suspicious bool) {
	trimmedMessage := strings.TrimSpace(message.Message[strings.Index(message.Message, " "):])

	// Check playlist
	playlistID := extractPlaylistId(trimmedMessage)
	if playlistID != "" {
		// TODO: analyze playlist
		reasonPlaylist, suspiciousPlaylist := analyzePlaylist(playlistID)
		suspicious = suspicious || suspiciousPlaylist
		appendStringComma(&reason, reasonPlaylist)
	}

	// Check video
	videoID := extractVideoId(trimmedMessage)
	log.Printf("Intercept | Check Video | %s | %s | \"%s\" | found ID \"%s\"\n",
		message.Channel, message.User.DisplayName, message.Message, videoID)

	reasonVideo, suspiciousVideo := analyzeVideo(videoID)
	suspicious = suspicious || suspiciousVideo
	appendStringComma(&reason, reasonVideo)

	return
}

// extractPlaylistId returns the playlist id if the message contains it ("&list=abc"), otherwise an empty string
func extractPlaylistId(message string) string {
	var playlistIdFromUrlRegex = regexp.MustCompile(`(?m)list=([^&\s]+)`)
	playlistIdMatches := playlistIdFromUrlRegex.FindStringSubmatch(message)
	if playlistIdMatches != nil && len(playlistIdMatches) >= 2 {
		return playlistIdMatches[1]
	}

	return ""
}

// extractVideoId finds a youtube video id in a message
// This finds IDs in youtube links, youtube shortlinks (youtu.be) and standalone IDs.
// This may return the entire message if it only consists of one word
// ==> This may return garbage or "", it does not check if the result is a valid (or existing) id
// (Youtube does not specify an ID format other than "string")
// This function assumes the ID does not contain "&" or spaces
func extractVideoId(message string) string {
	var videoIdFromUrlParameterRegex = regexp.MustCompile(`(?m)v=([^&\s]+)`)
	var videoIdFromShortUrl = regexp.MustCompile(`(?m)youtu.be/([^&\s]+)`)
	urlParameterMatches := videoIdFromUrlParameterRegex.FindStringSubmatch(message)
	if urlParameterMatches != nil && len(urlParameterMatches) >= 2 {
		return urlParameterMatches[1]
	}

	shortUrlMatches := videoIdFromShortUrl.FindStringSubmatch(message)
	if shortUrlMatches != nil && len(shortUrlMatches) >= 2 {
		return shortUrlMatches[1]
	}

	if !strings.ContainsAny(message, " ") {
		return message
	}
	return ""
}

func analyzeVideo(videoID string) (reason string, suspicious bool) {
	suspicious = false
	reason = ""
	videoList, err := youtubeService.Videos.List("snippet,contentDetails,status,statistics").Id(videoID).Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Analyze | Search | Video | %s | %d videos found\n", videoID, len(videoList.Items))
	for _, video := range videoList.Items {

		// TODO: Title word blacklist
		if strings.Contains(video.Snippet.Title, "Troll") {
			appendStringComma(&reason, "Title contains blacklisted word \"Troll\"")
			suspicious = true
		}

		if video.ContentDetails.ContentRating != nil && video.ContentDetails.ContentRating.YtRating == "ytAgeRestricted" {
			appendStringComma(&reason, "Video is age restricted")
			suspicious = true
		}

		// TODO: Description analysis

		if float64(video.Statistics.LikeCount)/float64(video.Statistics.LikeCount+video.Statistics.DislikeCount) < 0.75 {
			// TODO: make rating threshold configurable
			appendStringComma(&reason, "Video has less than 75% positive rating")
			suspicious = true
		}

		if video.Statistics.ViewCount < 10000 {
			// TODO: make viewcount threshold configurable
			appendStringComma(&reason, "Video has less than 10000 views")
			suspicious = true
		}

		if video.Statistics.CommentCount < 50 {
			appendStringComma(&reason, "Video has less than 50 comments")
			suspicious = true
		}

		if video.Status.UploadStatus != "processed" {
			appendStringComma(&reason, "Upload status is "+video.Status.UploadStatus+" (not processed)")
			suspicious = true
		}

		if video.Status.PrivacyStatus == "private" {
			appendStringComma(&reason, "Video is private")
			suspicious = true
		}

		// TODO: Analyze playlists containing this video

	}
	return
}

func analyzePlaylist(playlistID string) (reason string, suspicious bool) {
	suspicious = false
	reason = ""
	playlistList, err := youtubeService.Playlists.List("snippet,status").Id(playlistID).Do()
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Analyze | Search | Playlist | %s | %d playlists found\n", playlistID, len(playlistList.Items))
	for _, playlist := range playlistList.Items {

		if strings.Contains(playlist.Snippet.Title, "Troll") {
			appendStringComma(&reason, "Playlist title contains blacklisted word \"Troll\"")
			suspicious = true
		}

		if strings.Contains(playlist.Snippet.Description, "Troll") {
			appendStringComma(&reason, "Playlist description contains blacklisted word \"Troll\"")
			suspicious = true
		}

		for _, tag := range playlist.Snippet.Tags {
			if strings.EqualFold(tag, "troll") {
				appendStringComma(&reason, "Playlist taglist contains blacklisted word \""+tag+"\"")
				suspicious = true
			}
		}

		if playlist.Status.PrivacyStatus == "private" {
			appendStringComma(&reason, "Playlist status is private")
			suspicious = true
		}
	}
	// TODO: Video analysis of containing videos
	// TODO: Deleted content
	return reason, suspicious
}

func appendStringComma(target *string, suffix string) {
	if *target == "" {
		*target = suffix
	} else {
		*target = *target + ", " + suffix
	}
}
