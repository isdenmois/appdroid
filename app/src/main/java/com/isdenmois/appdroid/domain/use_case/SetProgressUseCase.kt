package com.isdenmois.appdroid.domain.use_case

import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.model.Download
import javax.inject.Inject

class SetProgressUseCase @Inject constructor() {
    operator fun invoke(app: AppPackage, download: Download) {
        app.progress.value = when (download) {
            is Download.Progress -> download.percent
            is Download.Finished ->  -1
        }
    }
}
