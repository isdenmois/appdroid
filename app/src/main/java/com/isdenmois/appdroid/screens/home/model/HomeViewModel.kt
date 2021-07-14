package com.isdenmois.appdroid.screens.home.model

import androidx.compose.runtime.mutableStateOf
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.isdenmois.appdroid.entities.apk.model.ApkInstaller
import com.isdenmois.appdroid.entities.apppackage.model.AppPackagesRepository
import com.isdenmois.appdroid.shared.api.AppPackage
import com.isdenmois.appdroid.shared.api.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import kotlinx.coroutines.flow.collect
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class HomeViewModel @Inject constructor(
    private val appPackagesRepository: AppPackagesRepository,
    private val apkInstaller: ApkInstaller,
): ViewModel() {
    var appList = mutableStateOf<Resource<List<AppPackage>>>(Resource.loading(null))

    init {
        fetchList()
    }

    fun fetchList() {
        viewModelScope.launch {
            appPackagesRepository.fetchApps().collect {
                appList.value = it
            }
        }
    }

    fun installApp(app: AppPackage) {
        viewModelScope.launch {
            apkInstaller.downloadAndInstallApplication(app.appId)
        }
    }
}
