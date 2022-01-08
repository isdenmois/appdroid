package com.isdenmois.appdroid.presentational.applist.ui

import androidx.compose.foundation.layout.Spacer
import androidx.compose.foundation.layout.fillMaxWidth
import androidx.compose.foundation.layout.height
import androidx.compose.foundation.layout.padding
import androidx.compose.material.LinearProgressIndicator
import androidx.compose.material.MaterialTheme
import androidx.compose.material.ProgressIndicatorDefaults
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.unit.dp
import com.isdenmois.appdroid.domain.model.AppPackage

@Composable
fun AppItemProgress(item: AppPackage) {
    val progress = item.progress.value / 100f

    if (progress > 0) {
        LinearProgressIndicator(
            progress = progress,
            color = MaterialTheme.colors.primaryVariant,
            backgroundColor = MaterialTheme.colors.background,
            modifier = Modifier
                .fillMaxWidth()
                .padding(horizontal = 16.dp)
        )
    } else {
        Spacer(modifier = Modifier.height(ProgressIndicatorDefaults.StrokeWidth))
    }
}
