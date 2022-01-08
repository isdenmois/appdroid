package com.isdenmois.appdroid.presentational.applist

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
import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.use_case.GetAppListUseCase
import com.isdenmois.appdroid.domain.use_case.InstallApplicationUseCase
import com.isdenmois.appdroid.domain.model.Resource
import dagger.hilt.android.lifecycle.HiltViewModel
import dagger.hilt.android.qualifiers.ApplicationContext
import kotlinx.coroutines.launch
import javax.inject.Inject

@HiltViewModel
class AppListViewModel @Inject constructor(
    @ApplicationContext private val applicationContext: Context,
    private val getAppList: GetAppListUseCase,
    private val installApplication: InstallApplicationUseCase,
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
            getAppList().collect {
                appList.value = it
            }
        }
    }

    fun installAppPackage(app: AppPackage) {
        viewModelScope.launch {
            installApplication(app)
        }
    }

    private fun enableWifi() {
        val wifiManager = applicationContext.getSystemService(ComponentActivity.WIFI_SERVICE) as WifiManager

        if (!wifiManager.isWifiEnabled) {
            wifiManager.isWifiEnabled = true
        }
    }
}
