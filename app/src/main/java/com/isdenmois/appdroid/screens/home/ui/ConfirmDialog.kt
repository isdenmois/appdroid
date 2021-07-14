package com.isdenmois.appdroid.screens.home.ui

import androidx.compose.foundation.layout.padding
import androidx.compose.material.AlertDialog
import androidx.compose.material.Button
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.isdenmois.appdroid.shared.api.AppPackage

@Composable
fun ConfirmDownloadDialog(app: AppPackage, onDismiss: () -> Unit, onConfirm: () -> Unit) {
    AlertDialog(
        modifier = Modifier.padding(24.dp),
        onDismissRequest = onDismiss,
        title = {
            Text("Download and install the app?")
        },
        confirmButton = {
            Button(onClick = { onDismiss(); onConfirm() }) {
                Text("Confirm")
            }
        },
        dismissButton = {
            Button(onClick = onDismiss) {
                Text("Dismiss")
            }
        },
        text = {
            Text(app.appId)
        },
    )
}
