package com.isdenmois.appdroid.screens.home.ui

import androidx.compose.foundation.border
import androidx.compose.foundation.layout.padding
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.AlertDialog
import androidx.compose.material.MaterialTheme
import androidx.compose.material.OutlinedButton
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.unit.dp
import com.isdenmois.appdroid.shared.api.AppPackage

@Composable
fun ConfirmDownloadDialog(app: AppPackage, onDismiss: () -> Unit, onConfirm: () -> Unit) {
    AlertDialog(
        modifier = Modifier
            .padding(24.dp)
            .border(width = 1.dp, color = Color.Black, shape = RoundedCornerShape(4.dp)),
        onDismissRequest = onDismiss,
        title = {
            Text("Download and install the app?")
        },
        confirmButton = {
            OutlinedButton(onClick = { onDismiss(); onConfirm() }) {
                Text("Confirm", color = MaterialTheme.colors.onPrimary)
            }
        },
        dismissButton = {
            OutlinedButton(onClick = onDismiss) {
                Text("Dismiss", color = MaterialTheme.colors.onSecondary)
            }
        },
        text = {
            Text(app.appId)
        },
    )
}
