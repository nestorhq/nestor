exports.handler = async (event, context, callback) => {
  console.log('>lambda>js>snoop> ', JSON.stringify(event, null, '  '));
  callback();
}
