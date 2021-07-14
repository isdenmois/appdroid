package com.isdenmois.appdroid.entities.apppackage.model

import android.content.Context
import android.content.pm.PackageInfo
import com.isdenmois.appdroid.shared.api.ApiService
import com.isdenmois.appdroid.shared.api.AppPackage
import com.isdenmois.appdroid.shared.api.Resource
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
