package com.isdenmois.appdroid.shared.ui

import androidx.compose.material.Text
import androidx.compose.ui.test.hasText
import androidx.compose.ui.test.junit4.createComposeRule
import com.isdenmois.appdroid.shared.api.Resource
import com.isdenmois.appdroid.ui.theme.AppDroidTheme
import org.junit.Rule
import org.junit.Test

class ResourceLoadingTest {
    @get:Rule
    val composeTestRule = createComposeRule()

    private fun <T>setContent(resource: Resource<T>) {
        composeTestRule.setContent {
            AppDroidTheme {
                ResourceLoading(resource = resource) {
                    Text("Success!")
                }
            }
        }
    }

    @Test
    fun onLoading() {
        setContent(Resource.loading(null))

        composeTestRule.onNode(hasText("Success!")).assertDoesNotExist()
    }

    @Test
    fun onSuccess() {
        setContent(Resource.success(data = emptyList<String>()))

        composeTestRule.onNode(hasText("Success!")).assertExists()
    }

    @Test
    fun onError() {
        setContent(Resource.error(data = null, message = "Something bad occurred"))

        composeTestRule.onNode(hasText("Success!")).assertDoesNotExist()
    }
}
