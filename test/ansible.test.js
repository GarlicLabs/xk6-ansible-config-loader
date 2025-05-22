import ansibleConfigLoader from 'k6/x/ansible-config-loader';

export default function () {
  const config = ansibleConfigLoader.getConfig(__ENV.CONFIG_PATH);
  for (const group of config.group_config) {
    console.log(JSON.stringify(group));
  }
}