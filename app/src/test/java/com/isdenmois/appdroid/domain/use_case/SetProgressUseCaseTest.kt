package com.isdenmois.appdroid.domain.use_case

import org.junit.Test
import strikt.api.expectThat
import strikt.assertions.*

import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.model.Download
import java.io.File

class SetProgressUseCaseTest {
    private val setProgressUseCase = SetProgressUseCase()

    @Test
    fun setDownloadProgress() {
        val app = AppPackage(id = "test", appId = "com.test", type = "HZ", version = "1.0", versionName = "1.0")
        expectThat(app.progress.value).isEqualTo(-1)

        setProgressUseCase(app, Download.Progress(0))
        expectThat(app.progress.value).isEqualTo(0)

        setProgressUseCase(app, Download.Progress(50))
        expectThat(app.progress.value).isEqualTo(50)

        setProgressUseCase(app, Download.Progress(100))
        expectThat(app.progress.value).isEqualTo(100)
    }

    @Test
    fun setDownloadFinished() {
        val app = AppPackage(id = "test", appId = "com.test", type = "HZ", version = "1.0", versionName = "1.0")
        expectThat(app.progress.value).isEqualTo(-1)

        setProgressUseCase(app, Download.Progress(0))
        expectThat(app.progress.value).isEqualTo(0)

        setProgressUseCase(app, Download.Finished(File("/tmp")))
        expectThat(app.progress.value).isEqualTo(-1)
    }
}
