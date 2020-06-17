c = get_config()
c.NotebookApp.ip = '0.0.0.0'
c.NotebookApp.open_browser = False
c.NotebookApp.port = 8888
c.NotebookApp.notebook_dir = '/opt/python/jupyter'
c.LabApp.user_settings_dir = '/opt/python/jupyterlab/lab/user-settings'
c.LabApp.workspaces_dir = '/opt/python/jupyterlab/lab/workspaces'
c.NotebookApp.password = u'sha1:dca3b5cb4801:54e2fa13f212ff3b6ef69a23bf1d1f11c74afa2c'
