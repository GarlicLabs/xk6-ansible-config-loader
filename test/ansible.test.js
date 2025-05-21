import ansibleConfigLoader from 'k6/x/ansible-config-loader';

export default function () {
  ansibleConfigLoader.getConfig(__ENV.CONFIG_PATH);
}