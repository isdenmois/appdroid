package com.isdenmois.appdroid.ui.main.adapter

import android.graphics.Color
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.recyclerview.widget.RecyclerView
import com.isdenmois.appdroid.R
import com.isdenmois.appdroid.data.model.App
import kotlinx.android.synthetic.main.app_item.view.*

interface OnAppClickListener {
    fun onAppClick(app: App);
}

class MainAdapter(val listener: OnAppClickListener): RecyclerView.Adapter<MainAdapter.DataViewHolder>() {
    private var apps: Array<App> = emptyArray()

    fun setApps(apps: Array<App>) {
        this.apps = apps
        notifyDataSetChanged()
    }

    override fun onCreateViewHolder(parent: ViewGroup, viewType: Int) =
        DataViewHolder(
            LayoutInflater.from(parent.context).inflate(R.layout.app_item, parent, false)
        )

    override fun getItemCount() = apps.size

    override fun onBindViewHolder(holder: DataViewHolder, position: Int) = holder.bind(apps[position])

    inner class DataViewHolder(itemView: View) : RecyclerView.ViewHolder(itemView) {
        fun bind(app: App) {
            itemView.tv_app_name.text = app.name ?: app.appId

            if (app.localVersion != null && app.localVersion != app.version) {
                itemView.tv_app_version.text = "${app.localVersionName} (${app.localVersion}) -> ${app.versionName} (${app.version})"
                itemView.tv_app_version.setTextColor(Color.RED);
            } else {
                itemView.tv_app_version.text = app.versionName
            }

            itemView.setOnClickListener {
                listener.onAppClick(app)
            }
        }
    }
}
