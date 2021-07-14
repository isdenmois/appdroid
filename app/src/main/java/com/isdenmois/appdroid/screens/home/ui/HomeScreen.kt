package com.isdenmois.appdroid.screens.home.ui

import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.runtime.*
import androidx.lifecycle.viewmodel.compose.viewModel
import com.isdenmois.appdroid.entities.apppackage.ui.AppItem
import com.isdenmois.appdroid.screens.home.model.HomeViewModel
import com.isdenmois.appdroid.shared.api.AppPackage
import com.isdenmois.appdroid.shared.api.Status
import com.isdenmois.appdroid.shared.ui.ResourceLoading

@Composable
fun HomeScreen(vm: HomeViewModel = viewModel()) {
    val appList by remember { vm.appList }
    var selectedApp: AppPackage? by remember { mutableStateOf(null) }

    Scaffold(
        topBar = {
            TopAppBar(title = { Text("AppDroid") }, backgroundColor = MaterialTheme.colors.primary, actions = {
                if (appList.status !== Status.LOADING) {
                    IconButton(onClick = { vm.fetchList() }) {
                        Icon(Icons.Filled.Refresh, null)
                    }
                }
            })
        },
    ) {
        ResourceLoading(resource = appList) {
            LazyColumn {
                items(appList.data ?: listOf()) { item ->
                    AppItem(item = item, onClick = { selectedApp = item })
                }
            }
        }
    }

    selectedApp?.let { app ->
        ConfirmDownloadDialog(
            app = app,
            onDismiss = { selectedApp = null },
            onConfirm = { vm.installApp(app) }
        )
    }
}
