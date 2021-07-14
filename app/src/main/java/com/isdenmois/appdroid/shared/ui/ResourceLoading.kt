package com.isdenmois.appdroid.shared.ui

import androidx.compose.foundation.layout.Arrangement
import androidx.compose.foundation.layout.Column
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.foundation.layout.padding
import androidx.compose.material.CircularProgressIndicator
import androidx.compose.material.MaterialTheme
import androidx.compose.material.Text
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.isdenmois.appdroid.shared.api.Resource
import com.isdenmois.appdroid.shared.api.Status

@Composable
fun <T>ResourceLoading(resource: Resource<T>, content: @Composable (data: T) -> Unit) {
    when (resource.status) {
        Status.LOADING -> {
            Column(modifier = Modifier.fillMaxSize(), horizontalAlignment = Alignment.CenterHorizontally, verticalArrangement = Arrangement.Center) {
                CircularProgressIndicator()
            }
        }
        Status.ERROR -> {
            Text(resource.message ?: "Error", color = MaterialTheme.colors.error, modifier = Modifier.padding(16.dp))
        }
        Status.SUCCESS -> {
            content(resource.data!!)
        }
    }
}
