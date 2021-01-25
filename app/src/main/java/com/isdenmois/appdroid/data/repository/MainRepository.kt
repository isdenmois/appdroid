package com.isdenmois.appdroid.data.repository

import com.isdenmois.appdroid.data.api.ApiService
import com.isdenmois.appdroid.data.api.OnApkDownloadedListener

class MainRepository(private val apiService: ApiService) {
    fun getApps() = apiService.getApps()
    fun getApk(appId: String, listener: OnApkDownloadedListener) = apiService.getApk(appId, listener)
    fun getPackage(appId: String) = apiService.getPackage(appId)
}
