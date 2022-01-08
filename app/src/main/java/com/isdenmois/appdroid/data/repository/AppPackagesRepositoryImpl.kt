package com.isdenmois.appdroid.data.repository

import android.content.Context
import android.content.pm.PackageInfo
import android.content.pm.PackageManager
import android.os.Environment
import com.isdenmois.appdroid.data.remote.ApiService
import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.model.Download
import com.isdenmois.appdroid.domain.repository.AppPackagesRepository
import com.isdenmois.appdroid.data.extenstions.downloadToFileWithProgress
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.Flow
import javax.inject.Inject

class AppPackagesRepositoryImpl @Inject constructor(
    @ApplicationContext private val context: Context,
    private val apiService: ApiService,
    private val packageManager: PackageManager,
): AppPackagesRepository {
    override suspend fun fetchApps() = apiService.getAppList().map { item ->
        val packageInfo = getPackage(item.appId)

        if (packageInfo != null) {
            item.localVersion = packageInfo.longVersionCode.toString()
            item.localVersionName = packageInfo.versionName
        }

        item
    }

    override suspend fun downloadApp(app: AppPackage): Flow<Download> {
        val fileName = "${app.appId}.apk"
        val dir = context.getExternalFilesDir(Environment.DIRECTORY_DOWNLOADS) ?: context.cacheDir

        return apiService.getApk(fileName).downloadToFileWithProgress(dir, fileName)
    }

    private fun getPackage(appId: String): PackageInfo? {
        val packages = packageManager.getInstalledPackages(0)

        return packages.firstOrNull { p -> p.packageName.equals(appId)}
    }
}
