package main

import (
	"math/rand"
	"strings"
)

// ShuffleStringSlice will randomize a string slice.
// I know that is a really bad shuffle logic (i won`t call this an algorithm,
// why? because i wrote and understand it :D)
// But this is YOUR chance to contribute to an open source project.
// Replace this by a cool one!
func ShuffleStringSlice(a []string) {
	for i := range a {
		j := rand.Intn(i + 1)
		a[i], a[j] = a[j], a[i]
	}
}

// crop is a modified "sub string" function allowing to limit a string length to a certain number of chars (from either start or end of string) and having a pre/postfix applied if the string really was cropped.
// content is the string to perform the operation on
// chars is the max number of chars of the string. Negative value means cropping from end of string.
// afterstring is the pre/postfix string to apply if cropping occurs.
// crop2space is true, then crop will be applied at nearest space. False otherwise.
//
// This function is a port from the TYPO3 CMS (written in PHP)
// @link https://github.com/TYPO3/TYPO3.CMS/blob/aae88a565bdbbb69032692f2d20da5f24d285cdc/typo3/sysext/frontend/Classes/ContentObject/ContentObjectRenderer.php#L4065
func Crop(content string, chars int, afterstring string, crop2space bool) string {
	if chars == 0 {
		return content
	}

	if len(content) < chars || (chars < 0 && len(content) < (chars*-1)) {
		return content
	}

	var cropedContent string
	truncatePosition := -1

	if chars < 0 {
		cropedContent = content[len(content)+chars:]
		if crop2space == true {
			truncatePosition = strings.Index(cropedContent, " ")
		}
		if truncatePosition >= 0 {
			cropedContent = cropedContent[truncatePosition+1:]
		}
		cropedContent = afterstring + cropedContent

	} else {
		cropedContent = content[:chars-1]
		if crop2space == true {
			truncatePosition = strings.LastIndex(cropedContent, " ")
		}
		if truncatePosition >= 0 {
			cropedContent = cropedContent[0:truncatePosition]
		}
		cropedContent += afterstring
	}

	return cropedContent
}
