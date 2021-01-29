package com.isdenmois.appdroid.ui.main.view

import android.content.Context
import android.os.Bundle
import android.view.Menu
import android.view.MenuItem
import android.view.View
import android.widget.Toast
import androidx.appcompat.app.AlertDialog
import androidx.appcompat.app.AppCompatActivity
import androidx.lifecycle.Observer
import androidx.recyclerview.widget.DividerItemDecoration
import androidx.recyclerview.widget.LinearLayoutManager
import com.isdenmois.appdroid.R
import com.isdenmois.appdroid.data.api.ApiService
import com.isdenmois.appdroid.data.api.OnApkDownloadedListener
import com.isdenmois.appdroid.data.model.App
import com.isdenmois.appdroid.data.repository.MainRepository
import com.isdenmois.appdroid.ui.main.adapter.MainAdapter
import com.isdenmois.appdroid.ui.main.adapter.OnAppClickListener
import com.isdenmois.appdroid.ui.main.viewmodel.MainViewModel
import com.isdenmois.appdroid.utils.ApkInstaller
import com.isdenmois.appdroid.utils.Status
import kotlinx.android.synthetic.main.activity_main.*
import java.io.File


class MainActivity: AppCompatActivity(), OnAppClickListener {
    private lateinit var mainRepository: MainRepository
    private lateinit var mainViewModel: MainViewModel
    private lateinit var adapter: MainAdapter

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContentView(R.layout.activity_main)
        setSupportActionBar(toolbar)
        setupServices()
        setupUI()
        setupObserver()
    }

    override fun onCreateOptionsMenu(menu: Menu): Boolean {
        menuInflater.inflate(R.menu.menu_main, menu)
        return true
    }

    override fun onOptionsItemSelected(item: MenuItem): Boolean {
        return when (item.itemId) {
            R.id.action_refresh -> {
                mainViewModel.fetchApps()
                true
            }

            else -> super.onOptionsItemSelected(item)
        }
    }

    private fun setupServices() {
        mainRepository = MainRepository(ApiService(this))
        mainViewModel = MainViewModel(mainRepository)
    }

    private fun setupUI() {
        rv_list.layoutManager = LinearLayoutManager(this)
        adapter = MainAdapter(this)
        rv_list.adapter = adapter
    }

    private fun setupObserver() {
        mainViewModel.getApps().observe(this, Observer {
            when (it.status) {
                Status.SUCCESS -> {
                    progressBar.visibility = View.GONE
                    adapter.setApps(it.data ?: emptyArray())
                    rv_list.visibility = View.VISIBLE
                }
                Status.LOADING -> {
                    progressBar.visibility = View.VISIBLE
                    rv_list.visibility = View.GONE
                }
                Status.ERROR -> {
                    rv_list.visibility = View.GONE
                    Toast.makeText(this, it.message, Toast.LENGTH_LONG).show()
                }
            }
        })
    }

    override fun onAppClick(app: App) {
        val ctx: Context = this

        AlertDialog.Builder(ctx)
            .setTitle(app.appId)
            .setMessage("Download and install the app?")
            .setIcon(android.R.drawable.ic_dialog_alert)
            .setPositiveButton(android.R.string.yes) { _, _ ->
                installApk(app)
            }
            .setNegativeButton(android.R.string.no, null)
            .show()
    }

    private fun installApk(app: App) {
        val ctx: Context = this

        mainRepository.getApk(app.appId, object : OnApkDownloadedListener {
            override fun onSuccess(file: File) {
                ApkInstaller.installApplication(ctx, file);
                mainViewModel.fetchApps()
            }
        })
    }
}
