package com.isdenmois.appdroid.domain.repository

import java.io.File

interface ApkRepository {
    fun installApplication(file: File)
}
