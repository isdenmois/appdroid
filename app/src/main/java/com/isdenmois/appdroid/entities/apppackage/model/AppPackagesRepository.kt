package com.isdenmois.appdroid.entities.apppackage.model

import android.content.Context
import android.content.pm.PackageInfo
import android.os.Environment
import com.isdenmois.appdroid.shared.api.*
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.Flow
import kotlinx.coroutines.flow.flow
import javax.inject.Inject

class AppPackagesRepository @Inject constructor(
    @ApplicationContext private val context: Context,
    private val apiService: ApiService
) {
    suspend fun fetchApps(): Flow<Resource<List<AppPackage>>> = flow {
        emit(Resource.loading(null))

        try {
            emit(Resource.success(fetchAppList()))
        } catch (exception: Exception) {
            exception.printStackTrace()
            emit(Resource.error(null, message = exception.message ?: "Error Occurred!"))
        }
    }

    suspend fun downloadApk(appId: String): Flow<Download> {
        val fileName = "$appId.apk"
        val dir = context.getExternalFilesDir(Environment.DIRECTORY_DOWNLOADS) ?: context.cacheDir

        return apiService.getApk(fileName).downloadToFileWithProgress(dir, "$appId.apk")
    }

    private suspend fun fetchAppList() = apiService.getAppList().map { item ->
        val packageInfo = getPackage(item.appId)

        if (packageInfo != null) {
            item.localVersion = packageInfo.longVersionCode.toString()
            item.localVersionName = packageInfo.versionName
        }

        item
    }

    private fun getPackage(appId: String): PackageInfo? {
        val packages = context.packageManager.getInstalledPackages(0)

        return packages.firstOrNull { p -> p.packageName.equals(appId)}
    }
}
