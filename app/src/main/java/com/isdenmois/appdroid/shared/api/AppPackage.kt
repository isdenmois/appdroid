package com.isdenmois.appdroid.shared.api

import androidx.compose.runtime.mutableStateOf

data class AppPackage(
    val id: String,
    var name: String?,
    val appId: String,
    val version: String,
    val versionName: String,
    val type: String,
    var localVersion: String? = null,
    var localVersionName: String? = null,
) {
    val progress = mutableStateOf(-1)
}
