package com.isdenmois.appdroid.shared.ui

import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.Row
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.padding
import androidx.compose.material.MaterialTheme
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.tooling.preview.Preview
import androidx.compose.ui.unit.dp
import com.isdenmois.appdroid.ui.theme.AppDroidTheme

@Composable
fun Item(title: String, subtitle: String? = null, error: String? = null, onClick: () -> Unit = {}) {
    Row(Modifier.clickable { onClick() }.fillMaxWidth().padding(16.dp, 8.dp)) {
       Column {
           Text(title, style = MaterialTheme.typography.body1, color = MaterialTheme.colors.onPrimary)

           if (subtitle !== null && error === null) {
               Text(subtitle, style = MaterialTheme.typography.caption, color = MaterialTheme.colors.onSecondary)
           }

           if (error !== null) {
               Text(error, style = MaterialTheme.typography.caption, color = MaterialTheme.colors.error)
           }
       }
    }
}

@Preview(showBackground = true)
@Composable
private fun ItemPreview() {
    AppDroidTheme {
        Column {
            Item(title = "Test")
            Item(title = "Bookly", subtitle = "4.16.0")
            Item(title = "Bookly", subtitle = "4.16.0", error = "4.16.0 -> 4.17.0")
        }
    }
}
