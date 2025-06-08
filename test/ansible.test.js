import ansibleConfigLoader from 'k6/x/ansible-config-loader';

export default function () {
  const config = ansibleConfigLoader.getConfig(__ENV.CONFIG_PATH);
  let allGroup = config.group_config.filter((g) => g.group_name === 'all');
  let k8sGroup = config.group_config.filter((g) => g.group_name === 'kubernetes');

  if (k8sGroup.length > 1 || k8sGroup.length === 0) {
    throw new Error('Expected exactly one group with name "kubernetes"');
  }
  if (k8sGroup[0].group_vars.length === 0) {
    throw new Error('Expected group "kubernetes" to have at least one group variable');
  }
  if (k8sGroup[0].hosts.length === 0) {
    throw new Error('Expected group "kubernetes" to have at least one host');
  }
  if (allGroup.length !== 1) {
    throw new Error('Expected exactly one group with name "all"');
  }
}