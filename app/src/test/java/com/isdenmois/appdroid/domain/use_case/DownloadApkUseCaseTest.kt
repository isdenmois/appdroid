package com.isdenmois.appdroid.domain.use_case

import com.isdenmois.appdroid.domain.fake.AppPackagesRepositoryFake
import com.isdenmois.appdroid.domain.model.AppPackage
import com.isdenmois.appdroid.domain.model.Download
import kotlinx.coroutines.ExperimentalCoroutinesApi
import kotlinx.coroutines.test.runTest

import org.junit.Test
import strikt.api.expectThat
import strikt.assertions.isEqualTo
import java.io.File

@ExperimentalCoroutinesApi
class DownloadApkUseCaseTest {
    private val appPackagesRepository = AppPackagesRepositoryFake()
    private val setProgress = SetProgressUseCase()
    private val downloadApk = DownloadApkUseCase(
        appPackagesRepository = appPackagesRepository,
        setProgress = setProgress
    )

    @Test
    operator fun invoke() = runTest {
        val mockFile = File("/tmp/download.apk")
        appPackagesRepository.setDownloadList(
            Download.Progress(0),
            Download.Progress(50),
            Download.Progress(100),
            Download.Finished(mockFile),
        )
        val app = AppPackage(id = "test", appId = "com.test", type = "HZ", version = "1.0", versionName = "1.0")

        expectThat(downloadApk(app)).isEqualTo(mockFile)
    }
}
