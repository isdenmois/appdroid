package com.isdenmois.appdroid.shared.api

import androidx.compose.runtime.mutableStateOf
import com.squareup.moshi.Json
import com.squareup.moshi.JsonClass

@JsonClass(generateAdapter = true)
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
    @Json(ignore = true)
    val progress = mutableStateOf(-1)
}
