package com.isdenmois.appdroid.entities.apppackage.ui

import androidx.compose.foundation.layout.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.tooling.preview.Preview
import com.isdenmois.appdroid.shared.api.AppPackage
import com.isdenmois.appdroid.shared.ui.Item

@Composable
fun AppItem(item: AppPackage, onClick: () -> Unit = {}) {
    val title = item.name ?: item.appId
    val version = item.versionName
    var error: String? = null

    if (item.localVersion !== null && item.localVersion != item.version) {
        error = "${item.localVersionName} (${item.localVersion}) -> ${item.versionName} (${item.version})"
    }

    Column {
        Item(title = title, subtitle = version, error = error, onClick = onClick)
        AppItemProgress(item)
    }
}

@Preview(showBackground = true)
@Composable
private fun AppItemPreview() {
    Column {
        AppItem(item = AppPackage(id = "1", name = "AppDroid", appId = "com.isden.appdroid", version = "1", versionName = "1.0", "hz"))
        AppItem(item = AppPackage(id = "2", name = "Bookly", appId = "com.isden.bookly", version = "5", versionName = "1.2.30", "hz", localVersion = "4", localVersionName = "1.2.0"))
        AppItem(item = AppPackage(id = "3", name = null, appId = "com.y.z", version = "10", versionName = "3.5.0", "hz"))
    }
}
