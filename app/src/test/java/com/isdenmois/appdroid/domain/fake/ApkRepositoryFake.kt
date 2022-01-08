package com.isdenmois.appdroid.domain.fake

import com.isdenmois.appdroid.domain.repository.ApkRepository
import java.io.File

class ApkRepositoryFake: ApkRepository {
    var lastFile: File? = null

    override fun installApplication(file: File) {
        lastFile = file
    }
}
