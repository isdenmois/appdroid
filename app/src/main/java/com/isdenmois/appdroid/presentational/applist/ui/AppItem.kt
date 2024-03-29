package com.isdenmois.appdroid.presentational.applist.ui

import androidx.compose.foundation.layout.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.tooling.preview.Preview
import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.presentational.ui.Item

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
        AppItem(item = AppPackage(appId = "com.isden.appdroid", name = "AppDroid", version = "1", versionName = "1.0", "hz"))
        AppItem(item = AppPackage(appId = "com.isden.bookly", name = "Bookly", version = "5", versionName = "1.2.30", localVersion = "4", localVersionName = "1.2.0"))
        AppItem(item = AppPackage(appId = "com.y.z", name = null, version = "10", versionName = "3.5.0", "hz"))
    }
}
