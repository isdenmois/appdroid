package com.isdenmois.appdroid.presentational.theme

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material.MaterialTheme
import androidx.compose.material.darkColors
import androidx.compose.material.lightColors
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

private val DarkColorPalette = darkColors(
    primary = Red,
    primaryVariant = Red,
    onPrimary = Color.White,
    secondary = Color.Black,
    onSecondary = LightGrey,
    background = Color.Black,
)

private val LightColorPalette = lightColors(
    primary = Color.White,
    primaryVariant = Grey,
    onPrimary = Color.Black,
    secondary = Color.White,
    onSecondary = Color.DarkGray,
)

@Composable
fun AppDroidTheme(darkTheme: Boolean = isSystemInDarkTheme(), content: @Composable() () -> Unit) {
    val colors = if (darkTheme) {
        DarkColorPalette
    } else {
        LightColorPalette
    }

    MaterialTheme(
        colors = colors,
        typography = Typography,
        shapes = Shapes,
        content = content
    )
}