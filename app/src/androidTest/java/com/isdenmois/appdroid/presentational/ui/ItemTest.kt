package com.isdenmois.appdroid.presentational.ui

import androidx.compose.ui.test.hasText
import androidx.compose.ui.test.junit4.createComposeRule
import com.isdenmois.appdroid.presentational.theme.AppDroidTheme
import org.junit.Rule
import org.junit.Test

class ItemTest {
    private val titleText = "Title Text Test"
    private val subtitleText = "Subtitle!"
    private val errorText = "What's wrong with you?"

    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun itemTestView() {
        composeTestRule.setContent {
            AppDroidTheme {
                Item(title = titleText, subtitle = subtitleText)
            }
        }

        composeTestRule.onNode(hasText(titleText)).assertExists()
        composeTestRule.onNode(hasText(subtitleText)).assertExists()
        composeTestRule.onNode(hasText(errorText)).assertDoesNotExist()
    }

    @Test
    fun errorVisibility() {
        composeTestRule.setContent {
            AppDroidTheme {
                Item(title = titleText, subtitle = subtitleText, error = errorText)
            }
        }

        composeTestRule.onNode(hasText(titleText)).assertExists()
        composeTestRule.onNode(hasText(subtitleText)).assertDoesNotExist()
        composeTestRule.onNode(hasText(errorText)).assertExists()
    }
}
