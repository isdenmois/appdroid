package com.isdenmois.appdroid.screens.home.model

import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.content.IntentFilter
import android.net.ConnectivityManager
import android.net.wifi.WifiManager
import androidx.activity.ComponentActivity
import androidx.compose.runtime.mutableStateOf
import androidx.lifecycle.ViewModel
import androidx.lifecycle.viewModelScope
import com.isdenmois.appdroid.entities.apk.model.ApkInstaller
import com.isdenmois.appdroid.entities.apppackage.model.AppPackagesRepository
import com.isdenmois.appdroid.shared.api.AppPackage
import com.isdenmois.appdroid.shared.api.Download
import com.isdenmois.appdroid.shared.api.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.flow.*
import kotlinx.coroutines.launch
import java.io.File
import javax.inject.Inject

@HiltViewModel
class HomeViewModel @Inject constructor(
    @ApplicationContext private val applicationContext: Context,
    private val appPackagesRepository: AppPackagesRepository,
    private val apkInstaller: ApkInstaller,
) : ViewModel() {
    var appList = mutableStateOf<Resource<List<AppPackage>>>(Resource.loading(null))

    private val receiver = object : BroadcastReceiver() {
        override fun onReceive(context: Context, intent: Intent) {
            fetchList()
        }
    }

    init {
        enableWifi()
        fetchList()

        val filter = IntentFilter().apply {
            addAction(WifiManager.WIFI_STATE_CHANGED_ACTION)
            addAction(ConnectivityManager.CONNECTIVITY_ACTION)
        }

        applicationContext.registerReceiver(receiver, filter)
    }

    override fun onCleared() {
        super.onCleared()
        applicationContext.unregisterReceiver(receiver)
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
            val apk = downloadApk(app)

            apk?.let {
                apkInstaller.installApplication(apk)
            }
        }
    }

    private suspend fun downloadApk(app: AppPackage): File? {
        val result = appPackagesRepository.downloadApk(app.appId).mapNotNull {
            when (it) {
                is Download.Progress -> app.progress.value = it.percent
                is Download.Finished -> {
                    app.progress.value = -1
                    return@mapNotNull it.file
                }
            }

            return@mapNotNull null
        }.lastOrNull()

        return result
    }

    private fun enableWifi() {
        val wifiManager = applicationContext.getSystemService(ComponentActivity.WIFI_SERVICE) as WifiManager

        if (!wifiManager.isWifiEnabled) {
            wifiManager.isWifiEnabled = true
        }
    }
}
