package com.isdenmois.appdroid.presentational.applist

import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.material.*
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.Refresh
import androidx.compose.runtime.*
import androidx.lifecycle.viewmodel.compose.viewModel
import com.isdenmois.appdroid.presentational.applist.ui.AppItem
import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.model.Status
import com.isdenmois.appdroid.presentational.ui.ResourceLoading

@Composable
fun AppListScreen(vm: AppListViewModel = viewModel()) {
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
            onConfirm = { vm.installAppPackage(app) }
        )
    }
}
