package com.isdenmois.appdroid.app

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import com.isdenmois.appdroid.screens.home.ui.HomeScreen
import com.isdenmois.appdroid.ui.theme.AppDroidTheme
import dagger.hilt.android.AndroidEntryPoint

@AndroidEntryPoint
class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            AppDroidTheme {
                HomeScreen()
            }
        }
    }
}
